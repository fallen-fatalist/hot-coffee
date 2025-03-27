package postgres

import (
	"database/sql"
	"fmt"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// Errors
var (
	ErrPeriodTypeInvalid = errors.New("incorrect period type provided")
	ErrIncorrectMenuItem = errors.New("incorrect menu item fetched, it's absent in menu item count struct")
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

func (r *orderRepository) Create(order entities.Order) (int64, error) {
	// Begin the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	insertOrderQuery := `
		INSERT INTO orders(customer_id, status)
		VALUES ($1, $2)
		RETURNING order_id
	`

	row := tx.QueryRow(insertOrderQuery, order.CustomerID, order.Status)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to inserting new order: %w", row.Err())
	}

	var orderID int64
	err = row.Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to scan order ID from row: %w", err)
	}

	// Insert order items
	insertOrderItemQuery := `
		INSERT INTO order_items(menu_item_id, order_id, quantity, customization_info)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(insertOrderItemQuery, item.ProductID, orderID, item.Quantity, item.CustomizationInfo)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	// Deduct order items
	err = inventoryRepositoryInstance.deductOrderItemsIngredients(tx, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to deduct order items ingredients: %w", err)
	}

	// Commit transaction if all succeeds
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to commit transaction after order creation: %w", err)
	}

	return orderID, nil
}

func (r *orderRepository) GetAll() ([]entities.Order, error) {
	query := `
	SELECT 	
		o.order_id, c.fullname, o.status, o.created_at,
		oi.menu_item_id, oi.quantity, oi.customization_info
	FROM
		orders o
	LEFT JOIN order_items oi USING(order_id)
	JOIN customers c USING(customer_id)
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []entities.Order
	var currentItem *entities.Order

	for rows.Next() {
		var (
			orderItemID       string
			customerID        string
			status            string
			createdAt         string
			menuItemIDString  sql.NullString
			quantity          sql.NullFloat64
			customizationInfo sql.NullString
		)

		if err := rows.Scan(&orderItemID, &customerID, &status, &createdAt, &menuItemIDString, &quantity, &customizationInfo); err != nil {
			return nil, err
		}

		if currentItem == nil || currentItem.ID != orderItemID {
			if currentItem != nil {
				orderItems = append(orderItems, *currentItem)
			}

			currentItem = &entities.Order{
				ID:           orderItemID,
				CustomerName: customerID,
				Items:        []entities.OrderItem{},
				Status:       status,
				CreatedAt:    createdAt,
			}
		}

		menuItemID, _ := strconv.Atoi(menuItemIDString.String)

		if menuItemIDString.Valid && quantity.Valid && customizationInfo.Valid {
			currentItem.Items = append(currentItem.Items, entities.OrderItem{
				ProductID:         menuItemID,
				Quantity:          int(quantity.Float64),
				CustomizationInfo: customizationInfo.String,
			})
		}
	}

	if currentItem != nil {
		orderItems = append(orderItems, *currentItem)
	}
	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderItems, nil
}

func (r *orderRepository) GetById(idStr string) (entities.Order, error) {
	// Parse the ID as an integer
	id, err := strconv.Atoi(idStr)
	var order entities.Order

	if err != nil {
		return order, ErrNonNumericID
	}

	query := `
	SELECT 	
		o.order_id, o.customer_id, o.status, o.created_at,
		oi.menu_item_id, oi.quantity, oi.customization_info
	FROM
		orders o
	LEFT JOIN
		order_items oi
	ON 
		o.order_id = oi.order_id
	WHERE o.order_id = $1
	`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	order.Items = []entities.OrderItem{}
	for rows.Next() {
		var (
			orderItemID       string
			customerID        string
			status            string
			createdAt         string
			menuItemID        sql.NullString
			quantity          sql.NullFloat64
			customizationInfo sql.NullString
		)

		if err := rows.Scan(&orderItemID, &customerID, &status, &createdAt, &menuItemID, &quantity, &customizationInfo); err != nil {
			return order, err
		}
		if order.ID == "" {
			order.ID = orderItemID
			order.CustomerName = customerID
			order.Status = status
			order.CreatedAt = createdAt
		}

		menuItemIDInteger, _ := strconv.Atoi(menuItemID.String)

		if menuItemID.Valid && quantity.Valid && customizationInfo.Valid {
			order.Items = append(order.Items, entities.OrderItem{
				ProductID:         menuItemIDInteger,
				Quantity:          int(quantity.Float64),
				CustomizationInfo: customizationInfo.String,
			})
		}
	}
	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return order, err
	}

	// Check if no rows were found
	if order.ID == "" {
		return order, sql.ErrNoRows
	}
	return order, nil
}

func (r *orderRepository) GetOrderRevenue(orderID int64) (totalOrderRevenue float64, err error) {
	// Common table expression query
	query := `
		WITH payment AS (
 	   		SELECT SUM(mi.price * oi.quantity) AS paymentSum
 	   		FROM order_items oi
 	   		JOIN menu_items mi USING(menu_item_id)
 	   		WHERE oi.order_id = $1
		),
		first_cost AS (
		    SELECT SUM(oi.quantity * mii.quantity * i.price) AS firstCost
		    FROM order_items oi 
		    JOIN menu_items_ingredients mii USING(menu_item_id)
		    JOIN inventory i USING(inventory_item_id)
		    WHERE oi.order_id = $1
		)
		SELECT p.paymentSum - fc.firstCost AS total_revenue
		FROM payment p, first_cost fc
	`

	err = r.db.QueryRow(query, orderID).Scan(&totalOrderRevenue)
	return totalOrderRevenue, err

}

func (r *orderRepository) Update(idStr string, order entities.Order) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}
	query := `
		UPDATE orders
		SET customer_id = $1, status = $2
		WHERE order_id = $3
	`
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, order.CustomerName, order.Status, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	deleteItemsQuery := `
		DELETE FROM order_items WHERE order_id = $1
	`
	_, err = tx.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	orderItemsQuery := `
		INSERT INTO order_items(menu_item_id, order_id, quantity, customization_info)
		VALUES ($1, $2, $3, $4)
	`
	for _, item := range order.Items {
		_, err = tx.Exec(orderItemsQuery, item.ProductID, id, item.Quantity, item.CustomizationInfo)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *orderRepository) Delete(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ErrNonNumericID
	}

	deleteOrderQuery := `
		DELETE FROM orders WHERE order_id = $1
	`
	_, err = r.db.Exec(deleteOrderQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetOrderedItemsCountByPeriod(period, month string, year int) (map[string]int, error) {
	var query string
	var args []interface{}
	var orderedItemsCount map[string]int = make(map[string]int)

	// If period is 'day'
	if period == "day" {
		query = `
			SELECT 
				EXTRACT(DAY FROM orders.created_at) AS day, 
				SUM(order_items.quantity) AS order_count 
			FROM 
				orders 
			JOIN 
				order_items USING(order_id) 
			WHERE 
				TO_CHAR(orders.created_at, 'FMMonth') = $1 
				AND EXTRACT(YEAR FROM orders.created_at) = $2
				AND orders.status = 'closed'
			GROUP BY 
				EXTRACT(DAY FROM orders.created_at)
			ORDER BY 
				day;
		`

		args = append(args, month, year)

	} else if period == "month" {
		query = `
			SELECT 
  		  	    TO_CHAR(orders.created_at, 'FMMonth') AS month, 
  		  	    SUM(order_items.quantity) AS order_count 
  		  	FROM 
  		  	    orders
			JOIN
				order_items USING(order_id) 
  		  	WHERE 
  		  	    EXTRACT(YEAR FROM orders.created_at) = $1
  		  	    AND orders.status = 'closed'
  		  	GROUP BY 
  		  	    TO_CHAR(orders.created_at, 'FMMonth'), EXTRACT(MONTH FROM orders.created_at)
  		  	ORDER BY 
  		  	    EXTRACT(MONTH FROM orders.created_at);
		`
		args = append(args, year)
	} else {
		return orderedItemsCount, ErrPeriodTypeInvalid
	}

	// Query the database
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return orderedItemsCount, err
	}
	defer rows.Close()

	// Scan the result into orderedItemsCount
	for rows.Next() {
		var key string
		var countRaw []byte // Use []byte to read the raw value from the database
		var count int
		if period == "day" {
			var day int
			if err := rows.Scan(&day, &countRaw); err != nil {
				return orderedItemsCount, err
			}
			key = fmt.Sprintf("%d", day)
		} else if period == "month" {
			if err := rows.Scan(&key, &countRaw); err != nil {
				return orderedItemsCount, err
			}
		}

		// Convert countRaw to string, parse as float, and convert to int
		countFloat, err := strconv.ParseFloat(string(countRaw), 64)
		if err != nil {
			return orderedItemsCount, fmt.Errorf("error parsing count: %v", err)
		}
		count = int(countFloat) // Convert to int (rounding down)

		orderedItemsCount[strings.ToLower(key)] = count
	}
	return orderedItemsCount, nil
}

func (r *orderRepository) GetOrderedMenuItemsCountByPeriod(startDate, endDate time.Time) (entities.OrderedMenuItemsCount, error) {
	menuItemsCount := entities.OrderedMenuItemsCount{}
	args := []interface{}{startDate, endDate}
	query := `
		SELECT 
            menu_items.name AS name,
            SUM(order_items.quantity) AS quantity
        FROM orders
        JOIN order_items USING(order_id)
        JOIN menu_items USING(menu_item_id)
        WHERE orders.created_at BETWEEN $1 AND $2
        GROUP BY menu_items.name
	`

	// Query the database
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return menuItemsCount, err
	}
	defer rows.Close()

	// Scan the result into orderedItemsCount
	for rows.Next() {
		var menuItemName string
		var countRaw []byte // Use []byte to read the raw value from the database
		var itemCount int

		if err := rows.Scan(&menuItemName, &countRaw); err != nil {
			return menuItemsCount, err
		}

		// Convert countRaw to string, parse as float, and convert to int
		countFloat, err := strconv.ParseFloat(string(countRaw), 64)
		if err != nil {
			return menuItemsCount, fmt.Errorf("error parsing count: %v", err)
		}
		itemCount = int(countFloat) // Convert to int (rounding down)

		switch menuItemName {
		case "Espresso":
			menuItemsCount.Espresso = itemCount
		case "Latte":
			menuItemsCount.Latte = itemCount
		case "Cappuccino":
			menuItemsCount.Cappuccino = itemCount
		case "Americano":
			menuItemsCount.Americano = itemCount
		case "Flat White":
			menuItemsCount.FlatWhite = itemCount
		case "Mocha":
			menuItemsCount.Mocha = itemCount
		case "Croissant":
			menuItemsCount.Croissant = itemCount
		case "Muffin":
			menuItemsCount.Muffin = itemCount
		case "Blueberry Muffin":
			menuItemsCount.BlueberryMuffin = itemCount
		case "Chocolate Chip Cookie":
			menuItemsCount.ChocolateChipCookie = itemCount
		case "Bagel":
			menuItemsCount.Bagel = itemCount
		case "Cheesecake":
			menuItemsCount.Cheesecake = itemCount
		case "Tiramisu":
			menuItemsCount.Tiramisu = itemCount
		case "Chocolate Cake":
			menuItemsCount.ChocolateCake = itemCount
		case "Vanilla Cupcake":
			menuItemsCount.VanillaCupcake = itemCount
		default:
			return menuItemsCount, ErrIncorrectMenuItem
		}
	}

	return menuItemsCount, nil
}

func (r *orderRepository) SetOrderStatusHistory(id int64, pastStatus, newStatus string) error {

	var query string
	var args []interface{}

	// Insert NULL for past_status if it's the first status change
	if pastStatus == "" {
		query = `
		INSERT INTO order_status_history(order_id, new_status)
		VALUES ($1, $2)
		`
		args = []interface{}{id, newStatus}
	} else {
		query = `
		INSERT INTO order_status_history(order_id, past_status, new_status)
		VALUES ($1, $2, $3)
		`
		args = []interface{}{id, pastStatus, newStatus}
	}

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetCustomerIDByName(fullname string, phone string) (int64, error) {
	var customer_id int64
	var query string
	var args []interface{}
	switch {
	case phone != "":
		query = ` 
		SELECT 
			customer_id
		FROM customers
		WHERE fullname = $1 AND phone = $2
		`
		args = []interface{}{fullname, phone}
	default:
		query = ` 
		SELECT 
			customer_id
		FROM customers
		WHERE fullname = $1
		`
		args = []interface{}{fullname}
	}

	row := r.db.QueryRow(query, args...)

	err := row.Scan(&customer_id)
	// Handle no existing customer
	if err == sql.ErrNoRows {
		insertQuery := `
			INSERT INTO customers 
				(fullname)
			VALUES 
				($1)
			RETURNING customer_idd
		`
		row := r.db.QueryRow(insertQuery, fullname, phone)
		if row.Err() != nil {
			return 0, err
		}

		err = row.Scan(&customer_id)
		if err != nil {
			return 0, err
		}

	} else if err != nil {
		return 0, err
	}
	return customer_id, nil
}

func (r *orderRepository) GetOrdersFullTextSearchReport(q string, minPrice, maxPrice int) ([]entities.OrderReport, error) {

	query := `
	SELECT 	
	o.order_id, c.fullname, array_agg(m.name) AS items, 
	ROUND(CAST(
		ts_rank(setweight(to_tsvector(c.fullname || ' ' || string_agg(m.name, ' ')), 'A'), 
		websearch_to_tsquery($1)) AS numeric), 2) AS relevance, 
	sum(m.price) AS total
	FROM orders o
	JOIN order_items oi USING(order_id)
	JOIN customers c USING(customer_id)
	JOIN menu_items m USING(menu_item_id)
	GROUP BY o.order_id, c.fullname
	HAVING to_tsvector(c.fullname || ' ' || string_agg(m.name, ' ')) @@ websearch_to_tsquery($1)
	ORDER BY relevance DESC;
	`

	rows, err := r.db.Query(query, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []entities.OrderReport{}
	for rows.Next() {
		var order entities.OrderReport
		var items []string

		err := rows.Scan(&order.ID, &order.CustomerName, pq.Array(&items), &order.Relevance, &order.Total)
		if err != nil {
			return nil, err
		}

		order.Items = items

		if (minPrice == 0 || order.Total >= float64(minPrice)) && (maxPrice == 0 || order.Total <= float64(maxPrice)) {
			order.Items = items
			orders = append(orders, order)
		}
	}

	return orders, nil
}
