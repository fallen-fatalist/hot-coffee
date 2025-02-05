package postgres

import (
	"database/sql"
	"errors"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"math"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

// Errors
var (
	ErrUnitNotExists = errors.New("the provided unit measure does not exists in database")
	ErrNonNumericID  = errors.New("non-numeric ID provided")
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
        INSERT INTO inventory (name, price, quantity, unit) 
        VALUES ($1, $2, $3, $4)
		`

	args := []interface{}{item.Name, item.Price, item.Quantity, item.Unit}

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
		if err := rows.Scan(&item.IngredientID, &item.Name, &item.Price, &item.Quantity, &item.Unit); err != nil {
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

	if err := row.Scan(&item.IngredientID, &item.Name, &item.Price, &item.Quantity, &item.Unit); err != nil {
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
			price = $3,
			quantity = $4, 
			unit = $5
		WHERE inventory_item_id = $1
		`

	args := []interface{}{id, item.Name, item.Price, item.Quantity, item.Unit}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil

}

func (r *inventoryRepository) SaveInventoryTransaction(idStr string, quantity float64) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	query := `
	INSERT INTO inventory_transactions(inventory_item_id, transaction_quantity)
	VALUES ($1, $2)
	`
	args := []interface{}{id, quantity}

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

func (r *inventoryRepository) GetPage(sortBy string, offset, rowCount int) (entities.PaginatedInventoryItems, error) {
	page := entities.PaginatedInventoryItems{
		Items: []entities.PageInventoryItem{},
	}

	// Get the total number of items for pagination
	var totalCount int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM inventory`).Scan(&totalCount)
	if err != nil {
		return page, err
	}

	query := `
		SELECT name, quantity, price
		FROM inventory
		ORDER BY $1 ASC
		LIMIT $3 OFFSET $2
	`

	// Calculate pagination details
	page.TotalPages = int(math.Ceil(float64(totalCount) / float64(rowCount)))
	page.CurrentPage = offset/rowCount + 1
	page.PageSize = rowCount
	page.HasNextPage = page.CurrentPage < page.TotalPages

	rows, err := r.db.Query(query, sortBy, offset, rowCount)
	if err != nil {
		return page, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		var item entities.PageInventoryItem
		if err := rows.Scan(&item.Name, &item.Price, &item.Quantity); err != nil {
			return page, err
		}
		page.Items = append(page.Items, item)
	}

	return page, nil
}

type inventoryCount struct {
	ingridientID    string
	ingridientCount float64
	itemCount       int
}

func (r *inventoryRepository) deductOrderItemsIngridients(tx *sql.Tx, orderID int64) error {
	// Part 1 Join Menu Items and their Ingridients
	// We fetch data to know how many ingridients to deduct

	// DO INDEXATION for query IDs
	joinQuery := `
		SELECT
			mii.inventory_item_id,
			mii.quantity as ingridient_quantity,
			oi.quantity as item_count
		FROM menu_items_ingridients mii
		JOIN order_items oi USING(menu_item_id)
		WHERE oi.order_id = $1
	`
	menuItemsIngridients := make([]entities.MenuItemIngredient, 0)

	rows, err := tx.Query(joinQuery, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var count inventoryCount
		if err := rows.Scan(&count.ingridientID, &count.ingridientCount, &count.itemCount); err != nil {
			tx.Rollback()
			return err
		}
		menuItemsIngridients = append(menuItemsIngridients, entities.MenuItemIngredient{
			IngredientID: count.ingridientID,
			Quantity:     count.ingridientCount * float64(count.itemCount),
		})
	}

	// Part 2 Deduction part
	deductQuery := `
        UPDATE inventory
		SET  
			quantity = quantity - $2 
		WHERE inventory_item_id = $1
	`

	for _, menuItemIngridient := range menuItemsIngridients {
		args := []interface{}{menuItemIngridient.IngredientID, menuItemIngridient.Quantity}
		_, err := tx.Exec(deductQuery, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}
