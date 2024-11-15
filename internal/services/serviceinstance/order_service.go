package serviceinstance

import (
	"errors"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// Errors
var (
	ErrEmptyCustomerName         = errors.New("empty customer name provided in order")
	ErrNegativeOrderItemQuantity = errors.New("negative order item quantity provided")
	ErrZeroOrderItemQuantity     = errors.New("zero order item quantity provided")
	ErrIncirrectOrderStatus      = errors.New("order status is incorrect")
	ErrCreatedAtField            = errors.New("created at field is not empty")
	ErrNoItemsInOrder            = errors.New("order has no items")
	ErrProductIsNotExist         = errors.New("product provided in order does not exist")
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

func (o *orderService) GetTotalSales() (int, error) {
	res := 0
	orders, err := o.GetOrders()
	if err != nil {
		return 0, err
	}

	for _, order := range orders {
		if order.Status == "closed" {
			res++
		}
	}
	return res, nil
}

func (o *orderService) GetPopularMenuItems() ([]entities.MenuItem, error) {
	items := make([]entities.MenuItem, 0)
	// orders, err := o.GetOrders()
	// if err != nil {
	// 	return nil, err
	// }
	// menuItems, err := MenuService.GetMenuItems()
	// if err != nil {
	// 	return nil, err
	// }

	// itemSales := make(map[int]int)
	// for _, order := range orders {
	// 	if order.Status == "closed" {
	// 		for _, item := range order.Items {
	// 			for _, menuItem := range menuItems {
	// 				if item.ProductID == menuItem.ID {
	// 					items = append(items, menuItem)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return items, nil
}
