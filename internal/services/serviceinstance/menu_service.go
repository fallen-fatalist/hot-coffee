package serviceinstance

import (
	"errors"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repositories/repository"
)

// Errors
var (
	ErrEmptyMenuItemID   = errors.New("empty menu item id provided")
	ErrEmptyMenuItemName = errors.New("empty menu item name provided")
	ErrNegativePrice     = errors.New("negative price in menu item provided")
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
		return ErrIDCollision
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
		item.Price = 0
	}
	return nil
}
