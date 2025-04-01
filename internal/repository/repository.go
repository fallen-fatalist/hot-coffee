package repository

import (
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/vo"
	"time"
)

type InventoryRepository interface {
	Create(item entities.InventoryItem) error
	GetAll() ([]entities.InventoryItem, error)
	GetById(id string) (entities.InventoryItem, error)
	Update(id string, item entities.InventoryItem) error
	Delete(id string) error
	// Pager for inventory items \\
	GetPage(sortBy string, offset, rowCount int) (entities.PaginatedInventoryItems, error)
}

type MenuRepository interface {
	Create(item entities.MenuItem) (int, error)
	AddPriceDifference(menu_item_id int, price_difference float64) error
	GetAll() ([]entities.MenuItem, error)
	GetById(id string) (entities.MenuItem, error)
	Update(id string, item entities.MenuItem) error
	Delete(id string) error
	GetMenusFullTextSearchReport(q string, minPrice, maxPrice int) ([]entities.MenuReport, error)
}

type OrderRepository interface {
	Create(order entities.Order) (int64, error)
	SetOrderStatusHistory(id int64, pastStatus, newStatus string) error
	GetAll() ([]entities.Order, error)
	GetById(id string) (entities.Order, error)
	GetOrderRevenue(id int64) (float64, error)
	Update(id string, order entities.Order) error
	Delete(id string) error
	GetOrderedItemsCountByPeriod(period, month string, year int) (map[string]int, error)
	GetOrderedMenuItemsCountByPeriod(startDate, endDate time.Time) (entities.OrderedMenuItemsCount, error)
	GetCustomerIDByName(fullname string, phone string) (int64, error)
	GetOrdersFullTextSearchReport(q string, minPrice, maxPrice int) ([]entities.OrderReport, error)
	FetchInventoryUpdates(orderIDs []int64) ([]vo.InventoryUpdate, error)
}

type Repository struct {
	Inventory InventoryRepository
	Menu      MenuRepository
	Order     OrderRepository
}
