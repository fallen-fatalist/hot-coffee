package cmd

import (
	"net/http"

	httpserver "hot-coffee/internal/infrastructure/server/http"
)

// Returns the router
func routes() http.Handler {
	mux := http.NewServeMux()

	// Orders:
	//     POST /orders: Create a new order.
	//     GET /orders: Retrieve all orders.
	mux.HandleFunc("/orders", httpserver.HandleOrders)
	// 	   POST /orders/open: Retrieve all open orders
	mux.HandleFunc("/orders/open", httpserver.HandleOpenOrders)

	//     GET /orders/{id}: Retrieve a specific order by ID.
	//     PUT /orders/{id}: Update an existing order.
	//     DELETE /orders/{id}: Delete an order.
	mux.HandleFunc("/orders/{id}", httpserver.HandleOrder)
	//     POST /orders/{id}/close: Close an order.
	mux.HandleFunc("/orders/{id}/close", httpserver.HandleOrderClose)
	//     POST /orders/{id}/in_progress: Start processing an order.
	mux.HandleFunc("/orders/{id}/in_progress", httpserver.HandleOrderInProgress)
	// GET /numberOfOrderedItems?startDate={startDate}&endDate={endDate}
	mux.HandleFunc("/orders/numberOfOrderedItems", httpserver.HandleNumberOfOrderedItems)

	// Inventory:
	//     POST /inventory: Add a new inventory item.
	//     GET /inventory: Retrieve all inventory items.
	mux.HandleFunc("/inventory", httpserver.HandleInventory)

	//     GET /inventory/{id}: Retrieve a specific inventory item.
	//     PUT /inventory/{id}: Update an inventory item.
	//     DELETE /inventory/{id}: Delete an inventory item.
	mux.HandleFunc("/inventory/{id}", httpserver.HandleInventoryItem)

	// Menu Items:
	//     POST /menu: Add a new menu item.
	//     GET /menu: Retrieve all menu items.
	mux.HandleFunc("/menu", httpserver.HandleMenu)

	//     GET /menu/{id}: Retrieve a specific menu item.
	//     PUT /menu/{id}: Update a menu item.
	//     DELETE /menu/{id}: Delete a menu item.
	mux.HandleFunc("/menu/{id}", httpserver.HandleMenuItem)

	// Aggregations:
	// GET /reports/total-sales: Get the total sales amount.
	mux.HandleFunc("/reports/total-sales", httpserver.HandleTotalSales)
	// GET /reports/popular-items: Get a list of popular menu items.
	mux.HandleFunc("/reports/popular-items", httpserver.HandlePopularItems)

	// New functionality
	// GET /reports/orderedItemsByPeriod?period={day|month}&month={month}
	mux.HandleFunc("/reports/orderedItemsByPeriod", httpserver.HandleOrderedItemsByPeriod)
	// GET /getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}
	mux.HandleFunc("/inventory/getLeftOvers", httpserver.HandleInventoryLeftovers)

	// GET /reports/search?q=chocolate cake&filter=menu,orders&minPrice=10&maxPrice=12
	mux.HandleFunc("/reports/search", httpserver.HandleFullTextSearchReport)
	// New functionality
	// GET /getLeftOvers?sortBy=quantity?page=1&pageSize=4

	// POST /orders/batch-process
	mux.HandleFunc("/orders/batch-process", httpserver.HandleBatchOrders)

	// Logging middleware applied
	middlewareAppliedMux := httpserver.RequestLoggingMiddleware(mux)
	middlewareAppliedMux = httpserver.HeadersMiddleware(middlewareAppliedMux)

	return middlewareAppliedMux
}
