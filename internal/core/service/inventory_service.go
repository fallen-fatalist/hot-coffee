package service

import "hot-coffee/internal/core/entity"

type InventoryService interface {
	CreateInventoryItem(item entity.InventoryItem) error
	GetInventoryItems() ([]entity.InventoryItem, error)
	GetInventoryItem(id string) (entity.InventoryItem, error)
	UpdateInventoryItem(id string, item entity.InventoryItem) error
	DeleteInventoryItem(id string) error
}
