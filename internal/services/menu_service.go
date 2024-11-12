package service

import "hot-coffee/internal/core/entities"

type MenuService interface {
	CreateMenuItem(item entities.MenuItem) error
	GetMenuItems() ([]entities.MenuItem, error)
	GetMenuItem(id string) (entities.MenuItem, error)
	UpdateMenuItem(id string, item entities.MenuItem) error
	DeleteMenuItem(id string) error
}
