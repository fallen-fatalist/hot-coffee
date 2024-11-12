package serviceinstance

import (
	"errors"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
	"log/slog"
	"os"
)

// errors
var (
	ErrNegativeQuantity = errors.New("negative quantity of inventory item")
	ErrIDCollision      = errors.New("id collision between id in request body and id in url")
)

type inventoryService struct {
	inventoryRepository repository.InventoryRepository
}

func NewInventoryService(storage repository.InventoryRepository) *inventoryService {
	if storage == nil {
		slog.Error("Error while creating Inventory service: Nil pointer repository provided")
		os.Exit(1)
	}
	return &inventoryService{storage}
}

func (s *inventoryService) CreateInventoryItem(item entities.InventoryItem) error {
	if item.Quantity < 0 {
		return ErrNegativeQuantity
	}
	return s.inventoryRepository.Create(item)
}

func (s *inventoryService) GetInventoryItems() ([]entities.InventoryItem, error) {
	return s.inventoryRepository.GetAll()
}

func (s *inventoryService) GetInventoryItem(id string) (entities.InventoryItem, error) {
	return s.inventoryRepository.GetById(id)
}

func (s *inventoryService) UpdateInventoryItem(id string, item entities.InventoryItem) error {
	if id != item.IngredientID {
		return ErrIDCollision
	} else if item.Quantity < 0 {
		return ErrNegativeQuantity
	}

	return s.inventoryRepository.Update(id, item)
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	return s.inventoryRepository.Delete(id)
}
