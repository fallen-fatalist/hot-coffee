package serviceinstance

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repository"
)

// errors
var (
	ErrNegativeInventoryItemQuantity = errors.New("negative quantity of inventory item")
	ErrEmptyInventoryItemID          = errors.New("empty id provided")
	ErrEmptyInventoryItemName        = errors.New("empty name provided")
	ErrEmptyUnit                     = errors.New("empty unit provided")
	ErrInventoryItemIDCollision      = errors.New("id collision between id in request body and id in url")
	ErrInventoryItemAlreadyExists    = errors.New("inventory item with such id already exists")
	ErrIngridientIDContainsSlash     = errors.New("ingridient id contains slash")
	ErrIngridientIDContainsSpace     = errors.New("ingridient id contains space")
	ErrInvalidSortValue              = errors.New("incorrect sort by value provided")
	ErrNegativePage                  = errors.New("negative page provided")
	ErrNegativePageSize              = errors.New("negative page size proovided")
	ErrInvalidInventoryPrice         = errors.New("invalid inventory item price provided")
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
	if err := validateInventoryItem(&item); err != nil {
		return err
	}
	return s.inventoryRepository.Create(item)
}

func (s *inventoryService) CreateInventoryTransaction(id string, quantity float64) error {
	return s.inventoryRepository.CreateInventoryTransaction(id, quantity)
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

func (s *inventoryService) GetLeftovers(sortBy string, page, pageSize int) (entities.PaginatedInventoryItems, error) {
	emptyPage := entities.PaginatedInventoryItems{}

	// Default values
	if sortBy == "" {
		sortBy = "inventory_item_id"
	}
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// Validation
	if sortBy != "quantity" && sortBy != "price" && sortBy != "inventory_item_id" {
		return emptyPage, ErrInvalidSortValue
	} else if page < 1 {
		return emptyPage, ErrNegativePage
	} else if pageSize < 1 {
		return emptyPage, ErrNegativePage
	}

	// Processing
	offset, rowCount := (page-1)*pageSize, pageSize

	return s.inventoryRepository.GetPage(sortBy, offset, rowCount)
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
	} else if strings.Contains(item.IngredientID, "/") {
		return ErrIngridientIDContainsSlash
	} else if strings.Contains(item.IngredientID, " ") {
		return ErrIngridientIDContainsSpace
	} else if item.Price <= 0 {
		return ErrInvalidInventoryPrice
	}

	if item.Quantity == 0 {
		item.Quantity = 0
	}

	return nil
}
