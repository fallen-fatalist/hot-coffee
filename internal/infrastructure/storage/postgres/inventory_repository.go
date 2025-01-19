package postgres

import (
	"database/sql"
	"errors"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Errors
var (
	ErrUnitNotExists = errors.New("the provided unit measure does not exists in database")
	ErrNonNumericID  = errors.New("Non-numeric ID provided")
)

type inventoryRepository struct {
	db *sql.DB
}

var inventoryRepositoryInstance *inventoryRepository

func NewInventoryRepository() *inventoryRepository {
	if inventoryRepositoryInstance != nil {
		return inventoryRepositoryInstance
	}

	db, err := openDB()
	if err != nil {
		slog.Error("Error while opening connection with PostgreSQL: ", "error:", err.Error())
		os.Exit(1)
	}

	inventoryRepositoryInstance = &inventoryRepository{
		db: db,
	}

	return inventoryRepositoryInstance
}

func (r *inventoryRepository) Create(item entities.InventoryItem) error {
	query := `
        INSERT INTO inventory (name, quantity, unit) 
        VALUES ($1, $2, $3)
		`

	args := []interface{}{item.Name, item.Quantity, item.Unit}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepository) GetAll() ([]entities.InventoryItem, error) {
	query := `
		SELECT * 
		FROM inventory
	`
	// Query to get multiple users
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	var items []entities.InventoryItem
	for rows.Next() {
		var item entities.InventoryItem
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *inventoryRepository) GetById(idStr string) (entities.InventoryItem, error) {
	id, err := strconv.Atoi(idStr)
	var item entities.InventoryItem

	if err != nil {
		return item, ErrNonNumericID
	}

	query := `
		SELECT * 
		FROM inventory
		WHERE inventory_item_id = $1
	`
	// Query to get multiple users
	row := r.db.QueryRow(query, id)

	if err := row.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit); err != nil {
		return item, err
	}

	return item, nil

}

func (r *inventoryRepository) Update(idStr string, item entities.InventoryItem) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	query := `
        UPDATE inventory
		SET 
			name = $2, 
			quantity = $3, 
			unit = $4
		WHERE inventory_item_id = $1
		`

	args := []interface{}{id, item.Name, item.Quantity, item.Unit}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (r *inventoryRepository) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	query := `
		DELETE FROM inventory
		WHERE inventory_item_id = $1
	`

	_, err = r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
