package cmd

import (
	controller "hot-coffee/internal/infrastructure/controller"
	"net/http"
)

// Returns the router
func routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Orders:
	//     POST /orders: Create a new order.
	//     GET /orders: Retrieve all orders.
	mux.HandleFunc("/orders", controller.HandleOrders)

	//     GET /orders/{id}: Retrieve a specific order by ID.
	//     PUT /orders/{id}: Update an existing order.
	//     DELETE /orders/{id}: Delete an order.
	//     POST /orders/{id}/close: Close an order.
	mux.HandleFunc("/orders/{id}", controller.HandleOrder)

	// Inventory:
	//     POST /inventory: Add a new inventory item.
	//     GET /inventory: Retrieve all inventory items.
	mux.HandleFunc("/inventory", controller.HandleInventory)

	//     GET /inventory/{id}: Retrieve a specific inventory item.
	//     PUT /inventory/{id}: Update an inventory item.
	//     DELETE /inventory/{id}: Delete an inventory item.
	mux.HandleFunc("/inventory/{id}", controller.HandleInventoryItem)

	// Menu Items:
	//     POST /menu: Add a new menu item.
	//     GET /menu: Retrieve all menu items.
	mux.HandleFunc("/menu", controller.HandleMenu)

	//     GET /menu/{id}: Retrieve a specific menu item.
	//     PUT /menu/{id}: Update a menu item.
	//     DELETE /menu/{id}: Delete a menu item.
	mux.HandleFunc("/menu/{id}", controller.HandleMenuItem)

	// Aggregations:
	//     GET /reports/total-sales: Get the total sales amount.
	//     GET /reports/popular-items: Get a list of popular menu items.
	mux.HandleFunc("/reports/total-sales", controller.HandleTotalSales)
	mux.HandleFunc("/reports/popular-items", controller.HandlePopularItems)

	return mux
}
