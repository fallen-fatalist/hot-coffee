package serviceinstance

import (
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

// Constants
const (
	OpenStatus       = "open"
	ClosedStatus     = "closed"
	InProgressStatus = "in progress"
)

var statuses = []string{OpenStatus, ClosedStatus, InProgressStatus}

// Errors
var (
	ErrEmptyCustomerName           = errors.New("empty customer name provided in order")
	ErrNegativeOrderItemQuantity   = errors.New("negative order item quantity provided")
	ErrZeroOrderItemQuantity       = errors.New("zero order item quantity provided")
	ErrIncorrectOrderStatus        = errors.New("order status is incorrect")
	ErrCreatedAtField              = errors.New("created at field is not empty")
	ErrNoItemsInOrder              = errors.New("order has no items")
	ErrProductIsNotExist           = errors.New("product provided in order does not exist in menu")
	ErrClosedOrderCannotBeModified = errors.New("closed order cannot be modified")
	ErrNotEnoughIgredient          = errors.New("not enough ingredients")
	// OrdersCountByPeriod errors
	ErrPeriodDayInvalid   = errors.New("incorrect period day provided")
	ErrPeriodTypeInvalid  = errors.New("incorrect period type provided")
	ErrPeriodMonthInvalid = errors.New("incorrect period month provided")
	// MenuItemsCountByPeriod
	ErrEndDateEarlierThanStartDate = errors.New("end date is earlier than start date")
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

func (s *orderService) CreateOrder(order entities.Order) error {
	order.Status = OpenStatus
	if err := validateOrder(order); err != nil {
		return err
	}

	// Fetch customer customer_id
	customerID, err := s.repository.GetCustomerIDByName(order.CustomerName, "")
	if err != nil {
		return fmt.Errorf("error while fetching the customer name: %w", err)
	}
	order.CustomerID = customerID

	orderID, err := s.repository.Create(order)
	if err != nil {
		return err
	}

	err = s.repository.SetOrderStatusHistory(orderID, "", order.Status)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Must be optimized in future, to reduce the number of database queries during the request execution
// Creates order concurrently
func (o *orderService) CreateOrders(orders []entities.Order) (vo.BatchResponse, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	processedOrders := []dto.ProcessedOrder{}

	for _, order := range orders {
		wg.Add(1)
		go func(order entities.Order) {
			defer wg.Done()

			var processedOrder dto.ProcessedOrder
			orderID, err := o.repository.Create(order)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if _, is := err.(*errors.ErrInsufficientIngredient); is {
					processedOrder.Reason = "insufficient inventory"
				} else {
					processedOrder.Reason = "not handled reason"
				}
				processedOrder.Status = "rejected"
			} else {
				processedOrder.Status = "accepted"
			}
			processedOrder.ID = orderID
			orderRevenue, err := o.repository.GetOrderRevenue(orderID)
			if err != nil {
				slog.Error("Error while calculating revenue for the created order: ", "order_id", orderID)
				processedOrder.Total = -1
				processedOrder.Reason = "Error occured while calculating the total revenue"
			} else {
				processedOrder.Total = orderRevenue
			}
			processedOrders = append(processedOrders, processedOrder)

		}(order)
	}
	wg.Wait()

	response := vo.BatchResponse{
		ProcessedOrders: processedOrders,
		Summary: vo.Summary{
			TotalOrders:      len(processedOrders),
			Accepted:         0,
			Rejected:         0,
			TotalRevenue:     0,
			InventoryUpdates: []vo.InventoryUpdate{},
		},
	}

	for _, order := range processedOrders {
		if order.Status == "accepted" {
			response.Summary.Accepted++
			response.Summary.TotalRevenue += order.Total
		} else {
			response.Summary.Rejected++
		}
	}

	return response, nil

}

func (s *orderService) GetOrders() ([]entities.Order, error) {
	return s.repository.GetAll()
}

func (s *orderService) GetOrder(id string) (entities.Order, error) {
	return s.repository.GetById(id)
}

func (s *orderService) GetOrderRevenue(orderIDstr string) (float64, error) {
	orderID, err := strconv.Atoi(orderIDstr)
	if err != nil {
		return 0, errors.NewErrNonIntegerID("order", orderIDstr)
	}

	return s.repository.GetOrderRevenue(int64(orderID))
}

func (s *orderService) UpdateOrder(idStr string, order entities.Order) error {
	orderDB, err := s.repository.GetById(idStr)
	if err != nil {
		return err
	}
	pastStatus := order.Status
	order.Status = orderDB.Status

	if orderDB.Status == "closed" || orderDB.Status == "in progress" {
		return ErrClosedOrderCannotBeModified
	}
	if err := validateOrder(order); err != nil {
		return err
	}
	if idStr != order.ID {
		return ErrInventoryItemIDCollision
	}

	// if orderDB.Status == "in progress" {
	// 	if err := validateSufficienceOfIngredients(order); err != nil {
	// 		return err
	// 	}
	// }

	// MUST DO: Add deduction logic

	if pastStatus != orderDB.Status {
		if err := s.updateOrderStatusHistory(idStr, pastStatus, order.Status); err != nil {
			return err
		}
	}

	return s.repository.Update(idStr, order)
}

func (s *orderService) DeleteOrder(id string) error {
	return s.repository.Delete(id)
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
	order.Status = ClosedStatus
	return s.repository.Update(id, order)

}

func validateOrder(order entities.Order) error {
	if order.CustomerName == "" {
		return ErrEmptyCustomerName
	} else if !utils.In(order.Status, statuses) {
		return ErrIncorrectOrderStatus
	} else if len(order.Items) == 0 {
		return ErrNoItemsInOrder
	}

	// Validate presence of order items in menu
	menuItems, err := MenuService.GetMenuItems()
	if err != nil {
		return err
	}
	menuItemsList := make(map[string]bool)
	for _, menuItem := range menuItems {
		menuItemsList[menuItem.ID] = true
	}

	// Products validation
	for _, item := range order.Items {
		if _, exists := menuItemsList[item.ProductID]; !exists {
			return ErrProductIsNotExist
		} else if item.Quantity < 0 {
			return ErrNegativeOrderItemQuantity
		} else if item.Quantity == 0 {
			return ErrZeroOrderItemQuantity
		}
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
		if order.Status == ClosedStatus {
			for _, orderItem := range order.Items {
				menuItem, err := MenuService.GetMenuItem(orderItem.ProductID)
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

	itemSalesCount := make(map[string]int)

	for _, order := range orders {
		if order.Status == ClosedStatus {
			for _, orderItem := range order.Items {
				itemSalesCount[orderItem.ProductID] += orderItem.Quantity
			}
		}
	}

	itemsSalesCount := make(entities.MenuItemSalesByCount, 0, len(itemSalesCount))
	for menuItemID, salesCount := range itemSalesCount {
		menuItem, err := MenuService.GetMenuItem(menuItemID)
		if err != nil {
			return nil, err
		}

		itemsSalesCount = append(itemsSalesCount, entities.MenuItemSales{
			ProductID:   menuItemID,
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
		if order.Status == OpenStatus {
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
	}

	if period == "day" && year == 0 {
		year = time.Now().Year()
	}

	var ok bool
	if month, ok = monthCapitalized[month]; !ok && month != "" {
		return orderedItemsCountByPeriod, ErrPeriodMonthInvalid
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

func (o *orderService) GetOrderedMenuItemsCountByPeriod(startDate, endDate time.Time) (entities.OrderedMenuItemsCount, error) {
	if diff := endDate.Sub(startDate); diff < 0 {
		return entities.OrderedMenuItemsCount{}, ErrEndDateEarlierThanStartDate
	}
	return o.repository.GetOrderedMenuItemsCountByPeriod(startDate, endDate)
}
