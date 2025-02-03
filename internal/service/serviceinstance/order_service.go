package serviceinstance

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"time"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repository"
)

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
	ErrNotEnoughIgridient          = errors.New("not enough ingridients")
	// OrdersCountByPeriod errors
	ErrPeriodDayInvalid   = errors.New("Incorrect period day provided")
	ErrPeriodTypeInvalid  = errors.New("Incorrect period type provided")
	ErrPeriodMonthInvalid = errors.New("Incorrect period month provided")
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
	order.Status = "open"
	if err := validateOrder(order); err != nil {
		return err
	}

	// if err := validateSufficienceOfIngridients(order); err != nil {
	// 	return err
	// }

	return s.repository.Create(order)
}

func (s *orderService) GetOrders() ([]entities.Order, error) {
	return s.repository.GetAll()
}

func (s *orderService) GetOrder(id string) (entities.Order, error) {
	return s.repository.GetById(id)
}

func (s *orderService) UpdateOrder(id string, order entities.Order) error {
	orderDB, err := s.repository.GetById(id)
	if err != nil {
		return err
	}
	order.Status = orderDB.Status

	if orderDB.Status == "closed" {
		return ErrClosedOrderCannotBeModified
	}
	if err := validateOrder(order); err != nil {
		return err
	}

	if err := validateSufficienceOfIngridients(order); err != nil {
		return err
	}

	return s.repository.Update(id, order)
}

func (s *orderService) DeleteOrder(id string) error {
	return s.repository.Delete(id)
}

func (s *orderService) CloseOrder(id string) error {
	order, err := s.GetOrder(id)
	if err != nil {
		return err
	}
	order.Status = "closed"
	return s.repository.Update(id, order)
}

func (s *orderService) SetInProgress(id string) error {
	order, err := s.GetOrder(id)
	if err != nil {
		return err
	}

	order.Status = "in progress"

	if err := validateSufficienceOfIngridients(order); err != nil {
		return err
	}

	return s.repository.Update(id, order)

}

func validateOrder(order entities.Order) error {
	if order.CustomerName == "" {
		return ErrEmptyCustomerName
	} else if order.Status != "open" && order.Status != "closed" {
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

func validateSufficienceOfIngridients(order entities.Order) error {
	ingridients := make(map[string]float64)
	for _, orderItem := range order.Items {
		menuItem, err := MenuService.GetMenuItem(orderItem.ProductID)
		if err != nil {
			return err
		}
		for _, ingridient := range menuItem.Ingredients {
			ingridients[ingridient.IngredientID] += ingridient.Quantity * float64(orderItem.Quantity)
		}
	}

	// Ingridients quantity check
	for ingridientID, quantity := range ingridients {
		inventoryItem, err := InventoryService.GetInventoryItem(ingridientID)
		if err != nil {
			return err
		}

		if inventoryItem.Quantity < quantity {
			return fmt.Errorf(ErrNotEnoughIgridient.Error()+": %s", ingridientID)
		}
	}

	// deduction after check
	if err := deductInventory(ingridients); err != nil {
		return err
	}

	// insert into inventory transaction
	if err := addInventoryTransactions(ingridients); err != nil {
		return err
	}

	return nil
}

func deductInventory(ingridientsCount map[string]float64) error {
	for ingridientID, quantity := range ingridientsCount {
		inventoryItem, err := InventoryService.GetInventoryItem(ingridientID)
		if err != nil {
			return err
		}
		inventoryItem.Quantity -= quantity
		if err := InventoryService.UpdateInventoryItem(ingridientID, inventoryItem); err != nil {
			return err
		}
	}
	return nil
}

func addInventoryTransactions(ingridientsCount map[string]float64) error {
	for ingridientID, quantity := range ingridientsCount {
		if err := InventoryService.SaveInventoryTransaction(ingridientID, quantity); err != nil {
			return err
		}
	}
	return nil
}

func (o *orderService) GetTotalSales() (entities.TotalSales, error) {
	var res float64 = 0.0
	orders, err := o.GetOrders()
	if err != nil {
		return entities.TotalSales{}, err
	}

	for _, order := range orders {
		if order.Status == "closed" {
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
		if order.Status == "closed" {
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
		if order.Status == "open" {
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
