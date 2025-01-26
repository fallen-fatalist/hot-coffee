package service

import (
	"hot-coffee/internal/core/entities"
	"time"
)

type InventoryService interface {
	CreateInventoryItem(item entities.InventoryItem) error
	SaveInventoryTransaction(id string, quantity float64) error
	GetInventoryItems() ([]entities.InventoryItem, error)
	GetInventoryItem(id string) (entities.InventoryItem, error)
	UpdateInventoryItem(id string, item entities.InventoryItem) error
	DeleteInventoryItem(id string) error
	GetLeftovers(sortBy string, page, pageSize int) (entities.PaginatedInventoryItems, error)
}

type MenuService interface {
	CreateMenuItem(item entities.MenuItem) error
	GetMenuItems() ([]entities.MenuItem, error)
	GetMenuItem(id string) (entities.MenuItem, error)
	UpdateMenuItem(id string, item entities.MenuItem) error
	DeleteMenuItem(id string) error
}

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
	GetOrderedItemsByPeriod(period, month string, year int) (entities.OrderedItemsCountByPeriod, error)
	GetOrderedMenuItemsCountByPeriod(startDate, endDate time.Time) (entities.OrderedMenuItemsCount, error)
}

type Service struct {
	InventoryService InventoryService
	MenuService      MenuService
	OrderService     OrderService
}
