package repository

import "hot-coffee/internal/core/entity"

type MenuRepository interface {
	Create(item entity.MenuItem) error
	GetAll() ([]entity.MenuItem, error)
	GetById(id string) (entity.MenuItem, error)
	Update(id string, item entity.MenuItem) error
	Delete(id string) error
}
