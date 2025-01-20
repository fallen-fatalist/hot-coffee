package repository

import "hot-coffee/internal/core/entities"

type InventoryRepository interface {
	Create(item entities.InventoryItem) error
	GetAll() ([]entities.InventoryItem, error)
	GetById(id string) (entities.InventoryItem, error)
	Update(id string, item entities.InventoryItem) error
	Delete(id string) error
	GetPage(sortBy string, offset, rowCount int) ([]entities.InventoryItem, error)
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
}

type Repository struct {
	Inventory InventoryRepository
	Menu      MenuRepository
	Order     OrderRepository
}
