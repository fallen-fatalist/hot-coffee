package serviceinstance

import (
	"errors"
	"fmt"
	"strings"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// Errors
var (
	ErrEmptyMenuItemID            = errors.New("empty menu item id provided")
	ErrEmptyMenuItemName          = errors.New("empty menu item name provided")
	ErrNegativePrice              = errors.New("negative price in menu item provided")
	ErrZeroPrice                  = errors.New("zero prince in menu item provided")
	ErrNegativeIngridientQuantity = errors.New("ingridient negative quantity provided")
	ErrZeroIngridientQuantity     = errors.New("ingridient quantity is zero")
	ErrIngridientDuplicate        = errors.New("duplicated ingridient provided in inventory item")
	ErrIngridientIsNotInInventory = errors.New("ingridient is not present in inventory")
	ErrMenuItemIDContainsSpace    = errors.New("menu item id contains space")
	ErrMenuItemIDContainsSlash    = errors.New("menu item id contains slash")
)

type menuService struct {
	menuRepository repository.MenuRepository
}

func NewMenuService(repository repository.MenuRepository) *menuService {
	if repository == nil {
		panic("nil repository provided")
	}
	return &menuService{repository}
}

func (s *menuService) CreateMenuItem(item entities.MenuItem) error {
	if err := validateMenuItem(item); err != nil {
		return err
	}

	return s.menuRepository.Create(item)
}

func (s *menuService) GetMenuItem(id string) (entities.MenuItem, error) {
	return s.menuRepository.GetById(id)
}

func (s *menuService) GetMenuItems() ([]entities.MenuItem, error) {
	return s.menuRepository.GetAll()
}

func (s *menuService) UpdateMenuItem(id string, item entities.MenuItem) error {
	if err := validateMenuItem(item); err != nil {
		return err
	}

	if id != item.ID {
		return ErrInventoryItemIDCollision
	}

	return s.menuRepository.Update(id, item)
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

	ingridientList := make(map[string]bool)
	inventoryIngridients := make(map[string]bool)

	// Fill the map
	inventoryItems, err := InventoryService.GetInventoryItems()
	if err != nil {
		return fmt.Errorf("error while getting inventory items: %s", err)
	}
	for _, inventoryItem := range inventoryItems {
		inventoryIngridients[inventoryItem.IngredientID] = true
	}

	// Ingridients validation
	for _, ingridient := range item.Ingredients {
		// Ingridient duplicate check
		if _, exists := ingridientList[ingridient.IngredientID]; exists {
			return ErrIngridientDuplicate
		} else {
			ingridientList[ingridient.IngredientID] = true
		}

		// Ingridient presence in inventory check
		if _, exists := inventoryIngridients[ingridient.IngredientID]; !exists {
			return ErrIngridientIsNotInInventory
		}

		// Quantity check
		if ingridient.Quantity < 0 {
			return ErrNegativeIngridientQuantity
		} else if ingridient.Quantity == 0 {
			return ErrZeroIngridientQuantity
		}
	}
	return nil
}
