package service

import "hot-coffee/internal/core/entities"

type InventoryService interface {
	CreateInventoryItem(item entities.InventoryItem) error
	GetInventoryItems() ([]entities.InventoryItem, error)
	GetInventoryItem(id string) (entities.InventoryItem, error)
	UpdateInventoryItem(id string, item entities.InventoryItem) error
	DeleteInventoryItem(id string) error
}
