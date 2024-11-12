package repositories

import "hot-coffee/internal/001Domain/entities"

type InventoryRepository interface {
	Create(item entities.InventoryItem) error
	GetAll() ([]entities.InventoryItem, error)
	GetById(id string) (entities.InventoryItem, error)
	Update(id string, item entities.InventoryItem) error
	Delete(id string) error
}
