package serviceinstance

import (
	"database/sql"
	"log/slog"
	"os"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"hot-coffee/internal/repository"
)

// errors
var (
	ErrNegativeInventoryItemQuantity = errors.New("negative quantity of inventory item")
	ErrEmptyInventoryItemID          = errors.New("empty id provided")
	ErrEmptyInventoryItemName        = errors.New("empty name provided")
	ErrInvalidUnit                   = errors.New("incorrect unit for inventory item provided")
	ErrInventoryItemIDCollision      = errors.New("id collision between id in request body and id in url")
	ErrInventoryItemAlreadyExists    = errors.New("inventory item with such id already exists")
	ErrIngredientIDContainsSlash     = errors.New("ingredient id contains slash")
	ErrIngredientIDContainsSpace     = errors.New("ingredient id contains space")
	ErrInvalidSortValue              = errors.New("incorrect sort by value provided")
	ErrNegativePage                  = errors.New("negative page provided")
	ErrNegativePageSize              = errors.New("negative page size proovided")
	ErrInvalidInventoryPrice         = errors.New("invalid inventory item price provided")
	ErrNonNumericIngredientID        = errors.New("non-integer ingredient id provided")
	ErrNegativeIngredientID          = errors.New("negative or zero ingredient id provided")
	ErrInventoryItemDoesntExist      = errors.New("inventory item with such id does not exist")
	ErrNoInventoryItems              = errors.New("no inventory items")
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
	if err := validateInventoryItem(&item); err != nil && err != ErrEmptyInventoryItemID {
		return err
	}

	if err := s.inventoryRepository.Create(item); err != nil {
		if errors.Is(err, errors.ErrIDAlreadyExists) {
			return ErrInventoryItemAlreadyExists
		}
	}
	return nil
}

// Not needed in service layer
// func (s *inventoryService) SaveInventoryTransaction(id string, quantity float64) error {
// 	return s.inventoryRepository.SaveInventoryTransaction(id, quantity)
// }

func (s *inventoryService) GetInventoryItems() ([]entities.InventoryItem, error) {
	items, err := s.inventoryRepository.GetAll()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoInventoryItems
		}
		return nil, err
	}
	return items, nil
}

func (s *inventoryService) GetInventoryItem(id string) (entities.InventoryItem, error) {
	item, err := s.inventoryRepository.GetById(id)

	if err == sql.ErrNoRows {
		return entities.InventoryItem{}, ErrInventoryItemDoesntExist
	}
	return item, err
}

func (s *inventoryService) UpdateInventoryItem(id string, item entities.InventoryItem) error {
	if err := validateInventoryItem(&item); err != nil {
		return err
	}

	if id != item.IngredientID {
		return ErrInventoryItemIDCollision
	}

	if err := s.inventoryRepository.Update(id, item); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInventoryItemDoesntExist
		}
		return err
	}
	return nil
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	if err := s.inventoryRepository.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInventoryItemDoesntExist
		}
		return err
	}
	return nil
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

// Validation for inventory items \\

var validUnits = map[string]bool{
	"grams":  true,
	"liters": true,
	"ml":     true,
	"kg":     true,
}

func isValidUnit(unit string) bool {
	return validUnits[unit]
}

func validateInventoryItem(item *entities.InventoryItem) error {
	// TODO: add inventory item id in repository
	// ID Validation
	err := isValidID(item.IngredientID)
	if errors.Is(err, ErrEmptyID) {
		return ErrEmptyInventoryItemID
	} else if errors.Is(err, ErrNonNumericID) {
		return ErrNonNumericIngredientID
	} else if errors.Is(err, ErrNegativeID) {
		return ErrNegativeIngredientID
	}

	// Other fields validation
	if item.Name == "" {
		return ErrEmptyInventoryItemName
	} else if !isValidUnit(item.Unit) {
		return ErrInvalidUnit
	} else if item.Quantity < 0 {
		return ErrNegativeInventoryItemQuantity
	} else if item.Price <= 0 {
		return ErrInvalidInventoryPrice
	}

	if item.Quantity == 0 {
		item.Quantity = 0
	}

	return nil
}
