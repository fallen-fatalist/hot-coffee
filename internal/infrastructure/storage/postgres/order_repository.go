package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"hot-coffee/internal/core/entities"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
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

func (r *orderRepository) Create(order entities.Order) error {
	query := `
		INSERT INTO orders(customer_id, status, created_at)
		VALUES ($1, $2, $3)
		RETURNING order_id
	`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var orderId int
	err = tx.QueryRow(query, order.CustomerName, order.Status, order.CreatedAt).Scan(&orderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	orderItemsQuery := `
		INSERT INTO order_items(menu_item_id, order_id, quantity, customization_info)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(orderItemsQuery, item.ProductID, orderId, item.Quantity, item.CustomizationInfo)
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

func (r *orderRepository) GetAll() ([]entities.Order, error) {
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
			menuItemID        sql.NullString
			quantity          sql.NullFloat64
			customizationInfo sql.NullString
		)

		if err := rows.Scan(&orderItemID, &customerID, &status, &createdAt, &menuItemID, &quantity, &customizationInfo); err != nil {
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

		if menuItemID.Valid && quantity.Valid && customizationInfo.Valid {
			currentItem.Items = append(currentItem.Items, entities.OrderItem{
				ProductID:         menuItemID.String,
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
		if menuItemID.Valid && quantity.Valid && customizationInfo.Valid {
			order.Items = append(order.Items, entities.OrderItem{
				ProductID:         menuItemID.String,
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
