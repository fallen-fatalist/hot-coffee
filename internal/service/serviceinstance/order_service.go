package serviceinstance

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"hot-coffee/internal/dto"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/utils"
	"hot-coffee/internal/vo"
)

// Errors
var (
	ErrEmptyCustomerName           = errors.New("empty customer name provided in order")
	ErrNegativeOrderItemQuantity   = errors.New("negative order item quantity provided")
	ErrZeroOrderItemQuantity       = errors.New("zero order item quantity provided")
	ErrIncorrectOrderStatus        = errors.New("order status is incorrect")
	ErrCreatedAtField              = errors.New("created at field is not empty")
	ErrNoItemsInOrder              = errors.New("order has no items")
	ErrClosedOrderCannotBeModified = errors.New("closed order cannot be modified")
	ErrNotEnoughIgredient          = errors.New("not enough ingredients")
	ErrNoOrders                    = errors.New("no orders in storage")
	ErrEmptyOrderID                = errors.New("empty order id provided")
	ErrNegativeOrderID             = errors.New("negative order id provided")
	ErrZeroOrderID                 = errors.New("zero order id provided")
	ErrNonNumericOrderID           = errors.New("non-numeric id provided")
	ErrOrderNotExists              = errors.New("order with such id does not exist")
	ErrOrderAlreadyExists          = errors.New("order with such id already exists")
	// OrdersCountByPeriod errors
	ErrPeriodDayInvalid   = errors.New("incorrect period day provided")
	ErrPeriodTypeInvalid  = errors.New("incorrect period type provided")
	ErrPeriodMonthInvalid = errors.New("incorrect period month provided")
	ErrParameterInvalid   = errors.New("inappropriate optional parameter")
	// MenuItemsCountByPeriod
	ErrEndDateEarlierThanStartDate = errors.New("end date is earlier than start date")
	ErrInvalidDate                 = errors.New("Invalid date for 'endDate' or 'startDate'. Expected format: DD-MM-YYYY.")
)

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) *orderService {
	if repository == nil {
		slog.Error("Error while creating Order service: Nil pointer repository provided")
		os.Exit(1)
	}
	return &orderService{repository}
}

func (s *orderService) CreateOrder(order entities.Order) (int64, error) {
	if order.Status == "" {
		order.Status = entities.OpenStatus
	}

	if err := validateOrder(&order); err != nil && err != ErrEmptyOrderID {
		return -1, err
	}

	// Fetch customer customer_id
	customerID, err := s.repository.GetCustomerIDByName(order.CustomerName, "")
	if err != nil {
		return -1, fmt.Errorf("error while fetching the customer name: %w", err)
	}
	order.CustomerID = customerID

	// TODO: Divide deduction from repository layer, move it to service layer
	orderID, err := s.repository.Create(order)
	if err != nil {
		if errors.Is(err, errors.ErrIDAlreadyExists) {
			return -1, ErrOrderAlreadyExists
		}
		return -1, fmt.Errorf("failed to create order in repository: %w", err)
	}

	// Set initial status of order
	err = s.repository.SetOrderStatusHistory(orderID, "", order.Status)
	if err != nil {
		return -1, fmt.Errorf("failed to save order status history: %w", err)
	}

	return orderID, nil
}

// TODO: Must be optimized in future, to reduce the number of database queries during the request execution
func (o *orderService) CreateOrders(orders []entities.Order) (vo.BatchResponse, error) {
	response := vo.BatchResponse{
		OrderReports: nil,
		Summary: vo.Summary{
			TotalOrders:      0,
			Accepted:         0,
			Rejected:         0,
			TotalRevenue:     0,
			InventoryUpdates: []vo.InventoryUpdate{},
		},
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	orderReports := []dto.OrderReport{}
	var orderIDs []int64 = []int64{}

	for _, order := range orders {
		wg.Add(1)

		go func(order entities.Order) {
			defer wg.Done()

			var orderReport dto.OrderReport
			orderID, err := o.CreateOrder(order)

			// Critical sectiton: changing "response" struct
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				var errInsufficientIngredient *errors.ErrInsufficientIngredient
				if errors.As(err, &errInsufficientIngredient) {
					orderReport.Reason = "insufficient inventory"
				} else if errors.Is(err, ErrEmptyCustomerName) {
					orderReport.Reason = "empty customer name"
				} else if errors.Is(err, ErrMenuItemNotExists) {
					orderReport.Reason = "non-existing menu item provided"
				} else if errors.Is(err, ErrNegativeOrderItemQuantity) {
					orderReport.Reason = "negative product quantity provided"
				} else if errors.Is(err, ErrZeroOrderItemQuantity) {
					orderReport.Reason = "zero product quantity provided"
				} else {
					orderReport.Reason = "failed to create order due to unhandled errors"
				}
				orderReport.Status = "rejected"
				response.Summary.Rejected++
				slog.Error("Error while creating the order: ", "error", err.Error())
			} else {
				order.ID = strconv.Itoa(int(orderID))
				orderReport.Status = "accepted"
				response.Summary.Accepted++
				orderRevenue, err := o.repository.GetOrderRevenue(orderID)
				if err != nil {
					slog.Error("Error while calculating revenue for the created order: ", "order_id", orderID)
					orderReport.Total = 0
					orderReport.Reason = "Error occured while calculating the total revenue"
				} else {
					orderReport.Total = orderRevenue
					response.Summary.TotalRevenue += orderRevenue
				}
				orderReport.ID = orderID
				orderIDs = append(orderIDs, orderID)
			}
			orderReport.CustomerName = order.CustomerName
			orderReports = append(orderReports, orderReport)
		}(order)

	}

	wg.Wait()

	if len(orderReports) != len(orders) {
		slog.Error("Incorrect number of orders processed:", "number of processed orders", len(orderReports), "number of orders provided to process", len(orders))
	}

	var err error
	response.OrderReports = orderReports
	response.Summary.TotalOrders = len(orderReports)
	response.Summary.InventoryUpdates, err = o.fetchInventoryUpdates(orderIDs)

	return response, err

}

// TODO: change the location from Order service to Inventory service

// Takes array of Order IDs and return the total inventory updates data
func (s *orderService) fetchInventoryUpdates(orderIDs []int64) (InventoryUpdates []vo.InventoryUpdate, err error) {
	// Validation
	for _, orderID := range orderIDs {
		if orderID < -1 {
			return nil, ErrNegativeOrderID
		}
	}
	return s.repository.FetchInventoryUpdates(orderIDs)
}

func (s *orderService) GetOrders() ([]entities.Order, error) {
	orders, err := s.repository.GetAll()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoOrders
		}
		return nil, err

	}
	return orders, nil
}

func (s *orderService) GetOrder(id string) (entities.Order, error) {
	err := isValidID(id)
	if err != nil {
		return entities.Order{}, err
	}

	order, err := s.repository.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Order{}, ErrOrderNotExists
		}
		return entities.Order{}, err
	}
	return order, nil
}

func (s *orderService) GetOrderRevenue(orderIDstr string) (float64, error) {
	orderID, err := strconv.Atoi(orderIDstr)
	if err != nil {
		return 0, errors.NewErrNonIntegerID("order", orderIDstr)
	}

	return s.repository.GetOrderRevenue(int64(orderID))
}

func (s *orderService) UpdateOrder(idStr string, order entities.Order) error {
	if err := validateOrder(&order); err != nil {
		return err
	}
	if idStr != order.ID {
		return ErrInventoryItemIDCollision
	}

	orderDB, err := s.repository.GetById(idStr)
	if err != nil {
		return err
	}
	pastStatus := orderDB.Status

	customerID, err := s.repository.GetCustomerIDByName(order.CustomerName, "")
	if err != nil {
		return err
	}

	if customerID != order.CustomerID {
		order.CustomerID = customerID
	}

	// TODO: Make atomic fetch and update
	err = s.repository.Update(idStr, order)
	if err != nil {
		return err
	}

	// if orderDB.Status == "in progress" {
	// 	if err := validateSufficienceOfIngredients(order); err != nil {
	// 		return err
	// 	}
	// }

	if pastStatus != order.Status {
		if err := s.updateOrderStatusHistory(idStr, pastStatus, order.Status); err != nil {
			return err
		}
	}

	return nil
}

func (s *orderService) DeleteOrder(id string) error {
	if err := s.repository.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderNotExists
		}
		return err
	}
	return nil
}

func (s *orderService) CloseOrder(idStr string) error {
	order, err := s.GetOrder(idStr)
	if err != nil {
		return err
	}

	if order.Status != "in progress" {
		return ErrIncorrectOrderStatus
	}

	// Update status and record the change
	if err := s.updateOrderStatusHistory(idStr, order.Status, "closed"); err != nil {
		return err
	}

	order.Status = "closed"
	return s.repository.Update(idStr, order)
}

func (s *orderService) updateOrderStatusHistory(idStr, oldStatus, newStatus string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return s.repository.SetOrderStatusHistory(int64(id), oldStatus, newStatus)
}

func (s *orderService) SetInProgress(id string) error {
	order, err := s.GetOrder(id)
	if err != nil {
		return err
	}
	order.Status = entities.ClosedStatus
	return s.repository.Update(id, order)

}

var menuItemsMap map[int]bool = make(map[int]bool)

func validateOrder(order *entities.Order) error {
	if order.CustomerName == "" {
		return ErrEmptyCustomerName
	} else if !utils.In(order.Status, entities.Statuses) {
		return ErrIncorrectOrderStatus
	} else if len(order.Items) == 0 {
		return ErrNoItemsInOrder
	}

	// Validate presence of order items in menu
	if len(menuItemsMap) == 0 {
		// TODO: On Update of menu items update this map
		menuItemsList, err := MenuService.GetMenuItems()
		if err != nil {
			return err
		}
		for _, menuItem := range menuItemsList {
			menuItemID, _ := strconv.Atoi(menuItem.ID)
			menuItemsMap[menuItemID] = true
		}
	}

	// Products validation
	for _, item := range order.Items {
		if _, exists := menuItemsMap[item.ProductID]; !exists {
			return ErrMenuItemNotExists
		} else if item.Quantity < 0 {
			return ErrNegativeOrderItemQuantity
		} else if item.Quantity == 0 {
			return ErrZeroOrderItemQuantity
		}
	}

	// ID Validation
	err := isValidID(order.ID)
	if errors.Is(err, ErrEmptyID) {
		return ErrEmptyOrderID
	} else if errors.Is(err, ErrNonNumericID) {
		return ErrNonNumericOrderID
	} else if errors.Is(err, ErrNegativeID) {
		return ErrNegativeOrderID
	} else if errors.Is(err, ErrZeroID) {
		return ErrZeroOrderID
	}

	return nil
}

// func validateSufficienceOfIngredients(order entities.Order) error {
// 	ingredients := make(map[string]float64)
// 	for _, orderItem := range order.Items {
// 		menuItem, err := MenuService.GetMenuItem(orderItem.MenuItemID)
// 		if err != nil {
// 			return err
// 		}
// 		for _, ingredient := range menuItem.Ingredients {
// 			ingredients[ingredient.IngredientID] += ingredient.Quantity * float64(orderItem.Quantity)
// 		}
// 	}

// 	// Ingredients quantity check
// 	for ingredientID, quantity := range ingredients {
// 		inventoryItem, err := InventoryService.GetInventoryItem(ingredientID)
// 		if err != nil {
// 			return err
// 		}

// 		if inventoryItem.Quantity < quantity {
// 			return fmt.Errorf(ErrNotEnoughIgridient.Error()+": %s", ingredientID)
// 		}
// 	}

// 	deduction after check
// 	if err := deductInventory(ingredients); err != nil {
// 		return err
// 	}

// 	return nil
// }

// TODO: MUST BE REPLACED
// func deductInventory(ingredientsCount map[string]float64) error {
// 	for ingredientID, quantity := range ingredientsCount {
// 		inventoryItem, err := InventoryService.GetInventoryItem(ingredientID)
// 		if err != nil {
// 			return err
// 		}
// 		inventoryItem.Quantity -= quantity
// 		if err := InventoryService.UpdateInventoryItem(ingredientID, inventoryItem); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Not needed in service layer
// func addInventoryTransactions(ingredientsCount map[string]float64) error {
// 	for ingredientID, quantity := range ingredientsCount {
// 		if err := InventoryService.SaveInventoryTransaction(ingredientID, quantity); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (o *orderService) GetTotalSales() (entities.TotalSales, error) {
	var res float64 = 0.0
	orders, err := o.GetOrders()
	if err != nil {
		return entities.TotalSales{}, err
	}

	for _, order := range orders {
		if order.Status == entities.ClosedStatus {
			for _, orderItem := range order.Items {
				productID := strconv.Itoa(orderItem.ProductID)
				menuItem, err := MenuService.GetMenuItem(productID)
				if err != nil {
					return entities.TotalSales{}, err
				}
				res += menuItem.Price * float64(orderItem.Quantity)
			}
		}
	}

	return entities.TotalSales{Total: res}, nil
}

// TODO: Refactor and optimize
func (o *orderService) GetPopularMenuItems() ([]entities.MenuItemSales, error) {
	orders, err := o.GetOrders()
	if err != nil {
		return nil, err
	}

	itemSalesCount := make(map[int]int)

	for _, order := range orders {
		if order.Status == entities.ClosedStatus {
			for _, orderItem := range order.Items {
				itemSalesCount[orderItem.ProductID] += orderItem.Quantity
			}
		}
	}

	itemsSalesCount := make(entities.MenuItemSalesByCount, 0, len(itemSalesCount))
	for menuItemID, salesCount := range itemSalesCount {
		menuItemIDString := strconv.Itoa(menuItemID)
		menuItem, err := MenuService.GetMenuItem(menuItemIDString)
		if err != nil {
			return nil, err
		}

		itemsSalesCount = append(itemsSalesCount, entities.MenuItemSales{
			ProductID:   menuItemIDString,
			ProductName: menuItem.Name,
			SalesCount:  salesCount,
		})
	}

	sort.Sort(itemsSalesCount)
	highestSales := make([]entities.MenuItemSales, 0)

	if len(itemSalesCount) > 0 {
		highestSales = append(highestSales, itemsSalesCount[0])
	}

	for idx := len(itemSalesCount) - 1; idx >= 1 && itemsSalesCount[idx] == itemsSalesCount[idx-1]; idx-- {
		highestSales = append(highestSales, itemsSalesCount[idx-1])
	}
	return highestSales, nil
}

func (o *orderService) GetOpenOrders() ([]entities.Order, error) {
	orders, err := o.repository.GetAll()
	if err != nil {
		return []entities.Order{}, nil
	}
	openOrders := []entities.Order{}
	for _, order := range orders {
		if order.Status == entities.OpenStatus {
			openOrders = append(openOrders, order)
		}
	}

	return openOrders, nil
}

var monthCapitalized = map[string]string{
	"january":   "January",
	"february":  "February", // Adjust for leap years as needed
	"march":     "March",
	"april":     "April",
	"may":       "May",
	"june":      "June",
	"july":      "July",
	"august":    "August",
	"september": "September",
	"october":   "October",
	"november":  "November",
	"december":  "December",
}

func (o *orderService) GetOrderedItemsByPeriod(period, month string, year int) (entities.OrderedItemsCountByPeriod, error) {
	orderedItemsCountByPeriod := entities.OrderedItemsCountByPeriod{}
	if period != "month" && period != "day" {
		return orderedItemsCountByPeriod, ErrPeriodTypeInvalid
	} else if period == "month" && year == 0 {
		return orderedItemsCountByPeriod, ErrPeriodMonthInvalid
	} else if period == "day" && month == "" {
		return orderedItemsCountByPeriod, ErrPeriodDayInvalid
	} else if year != 0 && month != "" {
		return orderedItemsCountByPeriod, ErrParameterInvalid
	}

	if period == "day" && year == 0 {
		year = time.Now().Year()
	}

	if month != "" {
		var ok bool
		if month, ok = monthCapitalized[month]; !ok {
			return orderedItemsCountByPeriod, ErrPeriodMonthInvalid
		}
	}

	itemsCount, err := o.repository.GetOrderedItemsCountByPeriod(period, month, year)
	if err != nil {
		return orderedItemsCountByPeriod, err
	}

	orderedItemsCountByPeriod.Period = period
	orderedItemsCountByPeriod.Month = strings.ToLower(month)
	orderedItemsCountByPeriod.Year = year
	orderedItemsCountByPeriod.OrderedItemsCount = itemsCount

	return orderedItemsCountByPeriod, nil
}

var dateLayout = "02.01.2006"

func (o *orderService) GetOrderedMenuItemsCountByPeriod(startDateStr, endDateStr string) (entities.OrderedMenuItemsCount, error) {
	// ✅ If both dates are empty, return all orders
	if startDateStr == "" && endDateStr == "" {
		return o.repository.GetOrderedMenuItemsCountByPeriod(time.Time{}, time.Time{})
	}
	// when no startDate or endDate
	if startDateStr == "" {
		endDate, err := time.Parse(dateLayout, endDateStr)
		if err != nil {
			return entities.OrderedMenuItemsCount{}, ErrInvalidDate
		}
		return o.repository.GetOrderedMenuItemsCountByPeriod(time.Time{}, endDate)
	} else if endDateStr == "" {
		startDate, err := time.Parse(dateLayout, startDateStr)
		if err != nil {
			return entities.OrderedMenuItemsCount{}, ErrInvalidDate
		}
		return o.repository.GetOrderedMenuItemsCountByPeriod(startDate, time.Time{})
	}

	//When both
	startDate, err := time.Parse(dateLayout, startDateStr)
	if err != nil {
		return entities.OrderedMenuItemsCount{}, ErrInvalidDate
	}

	endDate, err := time.Parse(dateLayout, endDateStr)
	if err != nil {
		return entities.OrderedMenuItemsCount{}, ErrInvalidDate
	}
	if diff := endDate.Sub(startDate); diff < 0 {
		return entities.OrderedMenuItemsCount{}, ErrEndDateEarlierThanStartDate
	}
	return o.repository.GetOrderedMenuItemsCountByPeriod(startDate, endDate)
}
