package serviceinstance

import (
	"errors"
	"sort"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// Errors
var (
	ErrEmptyCustomerName           = errors.New("empty customer name provided in order")
	ErrNegativeOrderItemQuantity   = errors.New("negative order item quantity provided")
	ErrZeroOrderItemQuantity       = errors.New("zero order item quantity provided")
	ErrIncirrectOrderStatus        = errors.New("order status is incorrect")
	ErrCreatedAtField              = errors.New("created at field is not empty")
	ErrNoItemsInOrder              = errors.New("order has no items")
	ErrProductIsNotExist           = errors.New("product provided in order does not exist")
	ErrClosedOrderCannotBeModified = errors.New("closed order cannot be modified")
)

type orderService struct {
	repository repository.OrderRepository
}

func NewOrderService(repository repository.OrderRepository) *orderService {
	if repository == nil {
		panic("empty repository provided to order service")
	}
	return &orderService{repository}
}

func (s *orderService) CreateOrder(order entities.Order) error {
	if order.Status == "" {
		order.Status = "open"
	}
	if err := validateOrder(order); err != nil {
		return err
	}
	return s.repository.Create(order)
}

func (s *orderService) GetOrders() ([]entities.Order, error) {
	return s.repository.GetAll()
}

func (s *orderService) GetOrder(id string) (entities.Order, error) {
	return s.repository.GetById(id)
}

func (s *orderService) UpdateOrder(id string, order entities.Order) error {
	if err := validateOrder(order); err != nil {
		return err
	}

	if order.Status == "closed" {
		return ErrClosedOrderCannotBeModified
	}

	return s.repository.Update(order.ID, order)
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

func validateOrder(order entities.Order) error {
	if order.CustomerName == "" {
		return ErrEmptyCustomerName
	} else if order.Status != "open" && order.Status != "closed" {
		return ErrIncirrectOrderStatus
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

func (o *orderService) GetTotalSales() (entities.TotalSales, error) {
	var res float64
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

	highestID := ""
	itemSalesCount := make(map[string]int)

	for _, order := range orders {
		if order.Status == "closed" {
			for _, orderItem := range order.Items {
				itemSalesCount[orderItem.ProductID] += orderItem.Quantity
				if itemSalesCount[orderItem.ProductID] > itemSalesCount[highestID] {
					highestID = orderItem.ProductID
				}
			}
		}
	}

	if err != nil {
		return nil, err
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
	if len(highestSales) > 0 {
		highestSales = append(highestSales, itemsSalesCount[len(itemsSalesCount)-1])
	}
	for idx := len(itemSalesCount) - 1; idx >= 1 && itemsSalesCount[idx] == itemsSalesCount[idx-1]; idx-- {
		highestSales = append(highestSales, itemsSalesCount[idx-1])
	}
	return highestSales, nil
}
