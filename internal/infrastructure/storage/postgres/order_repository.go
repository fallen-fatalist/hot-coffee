package postgres

import (
	"database/sql"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"os"
)

type orderRepository struct {
	db *sql.DB
}

var orderRepositoryInstance *orderRepository

func NewOrderRepository() *orderRepository {
	if orderRepositoryInstance != nil {
		return orderRepositoryInstance
	}

	db, err := openDB()
	if err != nil {
		slog.Error("Error while opening connection with PostgreSQL: ", "error:", err.Error())
		os.Exit(1)
	}

	orderRepositoryInstance = &orderRepository{
		db: db,
	}

	return orderRepositoryInstance
}

func (r *orderRepository) Create(order entities.Order) error {
	return nil
}

func (r *orderRepository) GetAll() ([]entities.Order, error) {
	return nil, nil
}

func (r *orderRepository) GetById(id string) (entities.Order, error) {
	return entities.Order{}, nil
}

func (r *orderRepository) Update(id string, order entities.Order) error {
	return nil
}

func (r *orderRepository) Delete(id string) error {
	return nil
}
