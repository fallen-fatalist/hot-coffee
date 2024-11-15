package repository

import "hot-coffee/internal/core/entities"

type OrderRepository interface {
	Create(order entities.Order) error
	GetAll() ([]entities.Order, error)
	GetById(id string) (entities.Order, error)
	Update(id string, order entities.Order) error
	Delete(id string) error
}
