package serviceinstance

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"hot-coffee/internal/core/entities"
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
	if err := validateMenuItem(item); err != nil {
		return err
	}

	id, err := s.menuRepository.Create(item)
	if err != nil {
		return err
	}

	return s.menuRepository.AddPriceDifference(id, int(item.Price))
}

func (s *menuService) GetMenuItem(id string) (entities.MenuItem, error) {
	return s.menuRepository.GetById(id)
}

func (s *menuService) GetMenuItems() ([]entities.MenuItem, error) {
	return s.menuRepository.GetAll()
}

func (s *menuService) UpdateMenuItem(idStr string, item entities.MenuItem) error {
	if err := validateMenuItem(item); err != nil {
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

func validateMenuItem(item entities.MenuItem) error {
	if item.ID == "" {
		return ErrEmptyMenuItemID
	} else if item.Name == "" {
		return ErrEmptyMenuItemName
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
	return nil
}
