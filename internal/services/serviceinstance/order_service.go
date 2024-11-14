package serviceinstance

import (
	"errors"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// Errors
var (
	ErrEmptyCustomerName         = errors.New("empty customer name provided in order")
	ErrEmptyOrderID              = errors.New("empty order id provided")
	ErrNegativeOrderItemQuantity = errors.New("negative order item quantity provided")
	ErrZeroOrderItemQuantity     = errors.New("zero order item quantity provided")
	ErrOrderIDCollision          = errors.New("order id collision between id in request body and id in url")
	
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
	} else if id != order.ID {
		return ErrOrderIDCollision
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
	} else if order.ID == "" {
		return ErrEmptyOrderID
	} else if order.Status != "open" && order.Status != "closed" {
		
	}

	for _, item := range order.Items {
		if item.Quantity < 0 {
			return ErrNegativeOrderItemQuantity
		} else if item.Quantity == 0 {
			return ErrZeroOrderItemQuantity
		} 
	}
	return nil
}
