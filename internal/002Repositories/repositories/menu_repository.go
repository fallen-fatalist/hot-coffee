package repositories

import "hot-coffee/internal/001Domain/entities"

type MenuRepository interface {
	Create(item entities.MenuItem) error
	GetAll() ([]entities.MenuItem, error)
	GetById(id string) (entities.MenuItem, error)
	Update(id string, item entities.MenuItem) error
	Delete(id string) error
}
