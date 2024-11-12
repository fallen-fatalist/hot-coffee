package service

import "hot-coffee/internal/core/entity"

type MenuService interface {
	CreateMenuItem(item entity.MenuItem) error
	GetMenuItems() ([]entity.MenuItem, error)
	GetMenuItem(id string) (entity.MenuItem, error)
	UpdateMenuItem(id string, item entity.MenuItem) error
	DeleteMenuItem(id string) error
	GetPopularMenuItems() []entity.MenuItem
}
