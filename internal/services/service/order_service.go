package service

import (
	"hot-coffee/internal/core/entities"
)

type OrderService interface {
	CreateOrder(order entities.Order) error
	GetOrders() ([]entities.Order, error)
	GetOrder(id string) (entities.Order, error)
	UpdateOrder(id string, order entities.Order) error
	DeleteOrder(id string) error
	CloseOrder(id string) error
	GetTotalSales() (entities.TotalSales, error)
	GetPopularMenuItems() ([]entities.MenuItemSales, error)
	GetOpenOrders() ([]entities.Order, error)
}
