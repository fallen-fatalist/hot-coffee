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

func (r *menuRepository) Create(item entities.MenuItem) (int, error) {
	// Query to insert a new menu item
	query := `
        INSERT INTO menu_items (name, description, price) 
        VALUES ($1, $2, $3)
        RETURNING menu_item_id
	`

	// Use a transaction to ensure atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return -1, err
	}

	// Insert the menu item and get the generated ID
	var menuItemID int
	err = tx.QueryRow(query, item.Name, item.Description, item.Price).Scan(&menuItemID)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Query to insert menu item ingredients
	ingredientQuery := `
        INSERT INTO menu_items_ingredients (menu_item_id, inventory_item_id, quantity)
        VALUES ($1, $2, $3)
	`

	// Insert each ingredient
	for _, ingredient := range item.Ingredients {
		_, err = tx.Exec(ingredientQuery, menuItemID, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return menuItemID, nil
}

func (r *menuRepository) GetAll() ([]entities.MenuItem, error) {
	query := `
		SELECT 
			mi.menu_item_id, mi.name, mi.description, mi.price, 
			mii.inventory_item_id, mii.quantity
		FROM 
			menu_items mi
		LEFT JOIN 
			menu_items_ingredients mii 
		ON 
			mi.menu_item_id = mii.menu_item_id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []entities.MenuItem
	// Temporary pointer for the current menu item
	var currentItem *entities.MenuItem

	for rows.Next() {
		var (
			menuItemID    string
			name          string
			description   string
			price         float64
			ingredientID  sql.NullString
			ingredientQty sql.NullFloat64
		)

		// Scan basic menu item fields and ingredient fields
		if err := rows.Scan(&menuItemID, &name, &description, &price, &ingredientID, &ingredientQty); err != nil {
			return nil, err
		}

		// Check if we're processing a new menu item
		if currentItem == nil || currentItem.ID != menuItemID {
			// Add the previous menu item to the slice (if any)
			if currentItem != nil {
				menuItems = append(menuItems, *currentItem)
			}

			// Create a new menu item
			currentItem = &entities.MenuItem{
				ID:          menuItemID,
				Name:        name,
				Description: description,
				Price:       price,
				Ingredients: []entities.MenuItemIngredient{},
			}
		}

		// If there is an ingredient, add it to the Ingredients slice
		if ingredientID.Valid && ingredientQty.Valid {
			currentItem.Ingredients = append(currentItem.Ingredients, entities.MenuItemIngredient{
				IngredientID: ingredientID.String,
				Quantity:     ingredientQty.Float64,
			})
		}
	}

	// Add the last menu item to the slice (if any)
	if currentItem != nil {
		menuItems = append(menuItems, *currentItem)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return menuItems, nil
}

func (r *menuRepository) GetById(idStr string) (entities.MenuItem, error) {
	// Parse the ID as an integer
	id, err := strconv.Atoi(idStr)
	var menuItem entities.MenuItem

	if err != nil {
		return menuItem, ErrNonNumericID
	}

	// Query to get menu item and its ingredients
	query := `
		SELECT 
			mi.menu_item_id, mi.name, mi.description, mi.price, 
			mii.inventory_item_id, mii.quantity
		FROM 
			menu_items mi
		LEFT JOIN 
			menu_items_ingredients mii 
		ON 
			mi.menu_item_id = mii.menu_item_id
		WHERE 
			mi.menu_item_id = $1
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return menuItem, err
	}
	defer rows.Close()

	// Initialize the menu item ingredients
	menuItem.Ingredients = []entities.MenuItemIngredient{}

	// Iterate over the rows
	for rows.Next() {
		var (
			menuItemID    string
			name          string
			description   string
			price         float64
			ingredientID  sql.NullString
			ingredientQty sql.NullFloat64
		)

		// Scan the row
		if err := rows.Scan(&menuItemID, &name, &description, &price, &ingredientID, &ingredientQty); err != nil {
			return menuItem, err
		}

		// Populate menuItem fields (only once)
		if menuItem.ID == "" {
			menuItem.ID = menuItemID
			menuItem.Name = name
			menuItem.Description = description
			menuItem.Price = price
		}

		// Append ingredients, if any
		if ingredientID.Valid && ingredientQty.Valid {
			menuItem.Ingredients = append(menuItem.Ingredients, entities.MenuItemIngredient{
				IngredientID: ingredientID.String,
				Quantity:     ingredientQty.Float64,
			})
		}
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return menuItem, err
	}

	// Check if no rows were found
	if menuItem.ID == "" {
		return menuItem, sql.ErrNoRows
	}

	return menuItem, nil
}

func (r *menuRepository) Update(idStr string, item entities.MenuItem) error {
	// Convert ID from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	// Use a transaction to ensure atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Update the main menu item
	query := `
        UPDATE menu_items
        SET 
            name = $2, 
            description = $3, 
            price = $4
        WHERE menu_item_id = $1
	`
	_, err = tx.Exec(query, id, item.Name, item.Description, item.Price)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing ingredients for the menu item
	deleteQuery := `
        DELETE FROM menu_items_ingredients 
        WHERE menu_item_id = $1
	`
	_, err = tx.Exec(deleteQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert updated ingredients
	insertQuery := `
        INSERT INTO menu_items_ingredients (menu_item_id, inventory_item_id, quantity)
        VALUES ($1, $2, $3)
	`
	for _, ingredient := range item.Ingredients {
		_, err = tx.Exec(insertQuery, id, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *menuRepository) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	deleteQuery := `
        DELETE FROM menu_items
        WHERE menu_item_id = $1
	`
	_, err = r.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepository) AddPriceDifference(id int, price_difference int) error {
	query := `
	INSERT INTO price_history(menu_item_id, price_difference)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(query, id, price_difference)
	if err != nil {
		return err
	}

	return nil
}

func (r *menuRepository) GetMenusFullTextSearchReport(q string, minPrice, maxPrice int) ([]entities.MenuReport, error) {
	query := `
	SELECT 
    m.menu_item_id, 
    m.name, 
    m.description, 
    m.price,
    ROUND(CAST(ts_rank(setweight(to_tsvector(m.name || ' ' || m.description), 'A'), 
    websearch_to_tsquery($1)) AS NUMERIC), 2) 
    AS relevance
	FROM menu_items m
	WHERE to_tsvector(m.name || ' ' || m.description) @@ websearch_to_tsquery($1)
	ORDER BY relevance DESC;
	`
	rows, err := r.db.Query(query, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	menus := []entities.MenuReport{}
	for rows.Next() {
		var menu entities.MenuReport

		err := rows.Scan(&menu.ID, &menu.Name, &menu.Description, &menu.Price, &menu.Relevance)
		if err != nil {
			return nil, err
		}

		if (minPrice == 0 || menu.Price >= float64(minPrice)) && (maxPrice == 0 || menu.Price <= float64(maxPrice)) {
			menus = append(menus, menu)
		}
	}
	return menus, nil
}
