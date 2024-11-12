package repository

import "hot-coffee/internal/core/entity"

type InventoryRepository interface {
	Create(item entity.InventoryItem) error
	GetAll() ([]entity.InventoryItem, error)
	GetById(id string) (entity.InventoryItem, error)
	Update(id string, item entity.InventoryItem) error
	Delete(id string) error
}
