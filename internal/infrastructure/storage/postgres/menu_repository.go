package postgres

import (
	"database/sql"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"os"
)

type menuRepository struct {
	db *sql.DB
}

var menuRepositoryInstance *menuRepository

func NewMenuRepository() *menuRepository {
	if menuRepositoryInstance != nil {
		return menuRepositoryInstance
	}

	db, err := openDB()
	if err != nil {
		slog.Error("Error while opening connection with PostgreSQL: ", "error:", err.Error())
		os.Exit(1)
	}

	menuRepositoryInstance = &menuRepository{
		db: db,
	}

	return menuRepositoryInstance
}

func (r *menuRepository) Create(item entities.MenuItem) error {
	return nil
}

func (r *menuRepository) GetAll() ([]entities.MenuItem, error) {
	return nil, nil
}

func (r *menuRepository) GetById(id string) (entities.MenuItem, error) {
	return entities.MenuItem{}, nil
}

func (r *menuRepository) Update(id string, order entities.MenuItem) error {
	return nil
}

func (r *menuRepository) Delete(id string) error {
	return nil
}
