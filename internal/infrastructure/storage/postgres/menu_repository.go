package postgres

import (
	"database/sql"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"os"
	"strconv"
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
	query := `
        INSERT INTO menu_items (name, description, price) 
        VALUES ($1, $2, $3)
		`

	args := []interface{}{item.Name, item.Description, item.Price}
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *menuRepository) GetAll() ([]entities.MenuItem, error) {
	query := `
		SELECT * 
		FROM menu_items
	`

	// Query to get menu items
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	var items []entities.MenuItem
	for rows.Next() {
		var item entities.MenuItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *menuRepository) GetById(idStr string) (entities.MenuItem, error) {
	id, err := strconv.Atoi(idStr)
	var item entities.MenuItem

	if err != nil {
		return item, ErrNonNumericID
	}

	query := `
		SELECT * 
		FROM menu_items
		WHERE menu_item_id = $1
	`
	// Query to get multiple users
	row := r.db.QueryRow(query, id)

	if err := row.Scan(&item.ID, &item.Name, &item.Description, &item.Price); err != nil {
		return item, err
	}

	return item, nil
}

func (r *menuRepository) Update(idStr string, item entities.MenuItem) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	query := `
        UPDATE menu_items
		SET 
			name = $2, 
			description = $3, 
			price = $4
		WHERE menu_item_id = $1
		`

	args := []interface{}{id, item.Name, item.Description, item.Price}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepository) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	query := `
		DELETE FROM menu_items
		WHERE menu_item_id = $1
	`

	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
