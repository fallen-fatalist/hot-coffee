package repository

import "hot-coffee/internal/core/entity"

type OrderRepository interface {
	Create(order entity.Order) error
	GetAll() ([]entity.Order, error)
	GetById(id string) (entity.Order, error)
	Update(id string, order entity.Order) error
	Delete(id string) error
	CloseOrder(id string) error
}
