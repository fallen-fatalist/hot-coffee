package serviceinstance

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"hot-coffee/internal/repository"
)

// Errors
var (
	ErrEmptyMenuItemID            = errors.New("empty menu item id provided")
	ErrEmptyMenuItemName          = errors.New("empty menu item name provided")
	ErrNegativePrice              = errors.New("negative price in menu item provided")
	ErrZeroPrice                  = errors.New("zero price in menu item provided")
	ErrNegativeIngredientQuantity = errors.New("ingredient negative quantity provided")
	ErrZeroIngredientQuantity     = errors.New("ingredient quantity is zero")
	ErrIngredientDuplicate        = errors.New("duplicated ingredient provided in inventory item")
	ErrIngredientIsNotInInventory = errors.New("ingredient is not present in inventory")
	ErrMenuItemIDContainsSpace    = errors.New("menu item id contains space")
	ErrMenuItemIDContainsSlash    = errors.New("menu item id contains slash")
	ErrNoMenuItems                = errors.New("no menu items")
	ErrMenuItemAlreadyExists      = errors.New("menu item with such id already exists")
	ErrMenuItemNotExists          = errors.New("menu item with such id not exists")
	ErrNegativeMenuID             = errors.New("negative menu item id provided")
	ErrZeroMenuID                 = errors.New("zero menu item id provided")
	ErrEmptyMenuItemDescription   = errors.New("menu item with empty description provided")
	ErrNegativeMenuItemID         = errors.New("negative menu item id provided")
	ErrNonNumericMenuItemID       = errors.New("non-numeric menu item id provided")
)

// TODO: Loading Menu items into memory
type menuService struct {
	menuRepository repository.MenuRepository
}

func NewMenuService(repository repository.MenuRepository) *menuService {
	if repository == nil {
		slog.Error("Error while creating Menu service: Nil pointer repository provided")
		os.Exit(1)
	}
	return &menuService{repository}
}

func (s *menuService) CreateMenuItem(item entities.MenuItem) error {
	if err := validateMenuItem(&item); err != nil && err != ErrEmptyMenuItemID {
		return err
	}

	id, err := s.menuRepository.Create(item)
	if err != nil {
		if errors.Is(err, errors.ErrIDAlreadyExists) {
			return ErrMenuItemAlreadyExists
		}
		return err
	}

	if err := s.menuRepository.AddPriceDifference(id, int(item.Price)); err != nil {
		return err
	}

	return nil
}

func (s *menuService) GetMenuItem(idStr string) (entities.MenuItem, error) {
	// ID validation
	if err := isValidID(idStr); err != nil {
		return entities.MenuItem{}, err
	}

	item, err := s.menuRepository.GetById(idStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.MenuItem{}, ErrMenuItemNotExists
		}
		return entities.MenuItem{}, err
	}
	return item, err
}

func (s *menuService) GetMenuItems() ([]entities.MenuItem, error) {
	items, err := s.menuRepository.GetAll()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoMenuItems
		}
		return nil, err
	}
	return items, err
}

func (s *menuService) UpdateMenuItem(idStr string, item entities.MenuItem) error {
	if err := validateMenuItem(&item); err != nil {
		return err
	}

	if idStr != item.ID {
		return ErrInventoryItemIDCollision
	}

	if err := s.menuRepository.Update(idStr, item); err != nil {
		return err
	}

	menuItem, err := s.GetMenuItem(idStr)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	return s.menuRepository.AddPriceDifference(id, int(menuItem.Price)-int(item.Price))
}

func (s *menuService) DeleteMenuItem(id string) error {
	return s.menuRepository.Delete(id)
}

func validateMenuItem(item *entities.MenuItem) error {
	if item.Name == "" {
		return ErrEmptyMenuItemName
	} else if item.Description == "" {
		return ErrEmptyMenuItemDescription
	} else if item.Price < 0 {
		return ErrNegativePrice
	} else if item.Price == 0 {
		return ErrZeroPrice
	} else if strings.Contains(item.ID, "/") {
		return ErrMenuItemIDContainsSlash
	} else if strings.Contains(item.ID, " ") {
		return ErrMenuItemIDContainsSpace
	}

	ingredientList := make(map[string]bool)
	inventoryIngredients := make(map[string]bool)

	// Fill the map
	inventoryItems, err := InventoryService.GetInventoryItems()
	if err != nil {
		return fmt.Errorf("error while getting inventory items: %s", err)
	}
	for _, inventoryItem := range inventoryItems {
		inventoryIngredients[inventoryItem.IngredientID] = true
	}

	// Ingredients validation
	for _, ingredient := range item.Ingredients {
		// Ingredient duplicate check
		if _, exists := ingredientList[ingredient.IngredientID]; exists {
			return ErrIngredientDuplicate
		} else {
			ingredientList[ingredient.IngredientID] = true
		}

		// Ingredient presence in inventory check
		if _, exists := inventoryIngredients[ingredient.IngredientID]; !exists {
			return ErrIngredientIsNotInInventory
		}

		// Quantity check
		if ingredient.Quantity < 0 {
			return ErrNegativeIngredientQuantity
		} else if ingredient.Quantity == 0 {
			return ErrZeroIngredientQuantity
		}
	}

	// ID Validation
	err = isValidID(item.ID)
	if errors.Is(err, ErrEmptyID) {
		return ErrEmptyMenuItemID
	} else if errors.Is(err, ErrNegativeID) {
		return ErrNegativeMenuID
	} else if errors.Is(err, ErrNonNumericID) {
		return ErrNonNumericMenuItemID
	} else if errors.Is(err, ErrZeroID) {
		return ErrZeroMenuID
	}
	return nil
}
