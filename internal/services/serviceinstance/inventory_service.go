package serviceinstance

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// errors
var (
	ErrNegativeInventoryItemQuantity = errors.New("negative quantity of inventory item")
	ErrEmptyInventoryItemID          = errors.New("empty id provided")
	ErrEmptyInventoryItemName        = errors.New("empty name provided")
	ErrEmptyUnit                     = errors.New("empty unit provided")
	ErrInventoryItemIDCollision      = errors.New("id collision between id in request body and id in url")
	ErrInventoryItemAlreadyExists    = errors.New("inventory item with such id already exists")
	ErrInventoryNameContainsSlash    = errors.New("inventory item name contains slash")
	ErrInventoryNameContainsSpace    = errors.New("inventory item name contains space")
)

type inventoryService struct {
	inventoryRepository repository.InventoryRepository
}

func NewInventoryService(storage repository.InventoryRepository) (*inventoryService, error) {
	if storage == nil {
		slog.Error("Error while creating Inventory service: Nil pointer repository provided")
		os.Exit(1)
	}

	return &inventoryService{storage}, nil
}

func (s *inventoryService) CreateInventoryItem(item entities.InventoryItem) error {
	if err := validateInventoryItem(&item); err != nil {
		return err
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
	if err := validateInventoryItem(&item); err != nil {
		return err
	}

	if id != item.IngredientID {
		return ErrInventoryItemIDCollision
	}

	return s.inventoryRepository.Update(id, item)
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	return s.inventoryRepository.Delete(id)
}

func validateInventoryItem(item *entities.InventoryItem) error {
	if item.IngredientID == "" {
		return ErrEmptyInventoryItemID
	} else if item.Name == "" {
		return ErrEmptyInventoryItemName
	} else if item.Unit == "" {
		return ErrEmptyUnit
	} else if item.Quantity < 0 {
		return ErrNegativeInventoryItemQuantity
	} else if strings.Contains(item.Name, "/") {
		return ErrInventoryNameContainsSlash
	} else if strings.Contains(item.Name, " ") {
		return ErrInventoryNameContainsSpace
	} else if item.Quantity == 0 {
		item.Quantity = 0
	}

	return nil
}
