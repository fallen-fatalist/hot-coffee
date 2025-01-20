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
)

// Errors
var (
	ErrPeriodTypeInvalid = errors.New("incorrect period type provided")
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

func (r *orderRepository) Update(id string, order entities.Order) error {
	return nil
}

func (r *orderRepository) Delete(id string) error {
	return nil
}

func (r *orderRepository) GetClosedOrdersCountByPeriod(period, month string, year int) (map[string]int, error) {
	var query string
	var args []interface{}
	var orderedItemsCount map[string]int = make(map[string]int)

	// If period is 'day'
	if period == "day" {
		query = `
			SELECT 
				EXTRACT(DAY FROM created_at) AS day, 
				COUNT(order_id) AS order_count 
			FROM 
				orders 
			WHERE 
				TO_CHAR(created_at, 'FMMonth') = $1 
				AND EXTRACT(YEAR FROM created_at) = $2
				AND status = 'completed'
			GROUP BY 
				EXTRACT(DAY FROM created_at)
			ORDER BY 
				day;
		`
		args = append(args, month, year)

	} else if period == "month" {
		query = `
			SELECT 
  		  	    TO_CHAR(created_at, 'FMMonth') AS month, 
  		  	    COUNT(order_id) AS order_count 
  		  	FROM 
  		  	    orders 
  		  	WHERE 
  		  	    EXTRACT(YEAR FROM created_at) = $1
  		  	    AND status = 'completed'
  		  	GROUP BY 
  		  	    TO_CHAR(created_at, 'FMMonth'), EXTRACT(MONTH FROM created_at)
  		  	ORDER BY 
  		  	    EXTRACT(MONTH FROM created_at);
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
		var count int
		if period == "day" {
			var day int
			if err := rows.Scan(&day, &count); err != nil {
				return orderedItemsCount, err
			}
			key = fmt.Sprintf("%d", day)
		} else if period == "month" {
			if err := rows.Scan(&key, &count); err != nil {
				return orderedItemsCount, err
			}
		}
		orderedItemsCount[strings.ToLower(key)] = count
	}
	return orderedItemsCount, nil
}
