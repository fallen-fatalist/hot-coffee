package repository

import (
	"hot-coffee/internal/core/entities"
	"time"
)

type InventoryRepository interface {
	Create(item entities.InventoryItem) error
	SaveInventoryTransaction(id string, quantity float64) error
	GetAll() ([]entities.InventoryItem, error)
	GetById(id string) (entities.InventoryItem, error)
	Update(id string, item entities.InventoryItem) error
	Delete(id string) error
	GetPage(sortBy string, offset, rowCount int) (entities.PaginatedInventoryItems, error)
}

type MenuRepository interface {
	Create(item entities.MenuItem) error
	GetAll() ([]entities.MenuItem, error)
	GetById(id string) (entities.MenuItem, error)
	Update(id string, item entities.MenuItem) error
	Delete(id string) error
}

type OrderRepository interface {
	Create(order entities.Order) error
	GetAll() ([]entities.Order, error)
	GetById(id string) (entities.Order, error)
	Update(id string, order entities.Order) error
	Delete(id string) error
	GetOrderedItemsCountByPeriod(period, month string, year int) (map[string]int, error)
	GetOrderedMenuItemsCountByPeriod(startDate, endDate time.Time) (entities.OrderedMenuItemsCount, error)
}

type Repository struct {
	Inventory InventoryRepository
	Menu      MenuRepository
	Order     OrderRepository
}
