package cmd

import (
	"net/http"

	"hot-coffee/internal/infrastructure/controllers"
)

// Returns the router
func routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Orders:
	//     POST /orders: Create a new order.
	//     GET /orders: Retrieve all orders.
	mux.HandleFunc("/orders", controllers.HandleOrders)

	//     GET /orders/{id}: Retrieve a specific order by ID.
	//     PUT /orders/{id}: Update an existing order.
	//     DELETE /orders/{id}: Delete an order.
	mux.HandleFunc("/orders/{id}", controllers.HandleOrder)
	//     POST /orders/{id}/close: Close an order.
	mux.HandleFunc("/orders/{id}/close", controllers.HandleOrderClose)

	// Inventory:
	//     POST /inventory: Add a new inventory item.
	//     GET /inventory: Retrieve all inventory items.
	mux.HandleFunc("/inventory", controllers.HandleInventory)

	//     GET /inventory/{id}: Retrieve a specific inventory item.
	//     PUT /inventory/{id}: Update an inventory item.
	//     DELETE /inventory/{id}: Delete an inventory item.
	mux.HandleFunc("/inventory/{id}", controllers.HandleInventoryItem)

	// Menu Items:
	//     POST /menu: Add a new menu item.
	//     GET /menu: Retrieve all menu items.
	mux.HandleFunc("/menu", controllers.HandleMenu)

	//     GET /menu/{id}: Retrieve a specific menu item.
	//     PUT /menu/{id}: Update a menu item.
	//     DELETE /menu/{id}: Delete a menu item.
	mux.HandleFunc("/menu/{id}", controllers.HandleMenuItem)

	// Aggregations:
	//     GET /reports/total-sales: Get the total sales amount.
	//     GET /reports/popular-items: Get a list of popular menu items.
	mux.HandleFunc("/reports/total-sales", controllers.HandleTotalSales)
	mux.HandleFunc("/reports/popular-items", controllers.HandlePopularItems)
	mux.HandleFunc("/reports/open", controllers.HandleOpenOrders)

	return mux
}
