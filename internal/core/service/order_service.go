package service

import (
	"hot-coffee/internal/core/entity"
)

type OrderService interface {
	CreateOrder(order entity.Order) error
	GetOrders() ([]entity.Order, error)
	GetOrder(id string) (entity.Order, error)
	UpdateOrder(id string, order entity.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
	GetTotalSales() (int, error)
}
