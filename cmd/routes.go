package cmd

import (
	"net/http"

	h "hot-coffee/internal/infrastructure/server/http"
)

// Returns the router
func routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Orders:
	//     POST /orders: Create a new order.
	//     GET /orders: Retrieve all orders.
	mux.HandleFunc("/orders", h.HandleOrders)
	// 	   POST /orders/open: Retrieve all open orders
	mux.HandleFunc("/orders/open", h.HandleOpenOrders)

	//     GET /orders/{id}: Retrieve a specific order by ID.
	//     PUT /orders/{id}: Update an existing order.
	//     DELETE /orders/{id}: Delete an order.
	mux.HandleFunc("/orders/{id}", h.HandleOrder)
	//     POST /orders/{id}/close: Close an order.
	mux.HandleFunc("/orders/{id}/close", h.HandleOrderClose)

	// Inventory:
	//     POST /inventory: Add a new inventory item.
	//     GET /inventory: Retrieve all inventory items.
	mux.HandleFunc("/inventory", h.HandleInventory)

	//     GET /inventory/{id}: Retrieve a specific inventory item.
	//     PUT /inventory/{id}: Update an inventory item.
	//     DELETE /inventory/{id}: Delete an inventory item.
	mux.HandleFunc("/inventory/{id}", h.HandleInventoryItem)

	// Menu Items:
	//     POST /menu: Add a new menu item.
	//     GET /menu: Retrieve all menu items.
	mux.HandleFunc("/menu", h.HandleMenu)

	//     GET /menu/{id}: Retrieve a specific menu item.
	//     PUT /menu/{id}: Update a menu item.
	//     DELETE /menu/{id}: Delete a menu item.
	mux.HandleFunc("/menu/{id}", h.HandleMenuItem)

	// Aggregations:
	//     GET /reports/total-sales: Get the total sales amount.
	//     GET /reports/popular-items: Get a list of popular menu items.
	mux.HandleFunc("/reports/total-sales", h.HandleTotalSales)
	mux.HandleFunc("/reports/popular-items", h.HandlePopularItems)

	return mux
}
