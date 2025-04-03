package flag

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Global flags
var (
	StoragePath = "data"
	Port        = 4000
)

func Parse(args []string) (err error) {
	for _, arg := range args {
		if arg == "--help" {
			PrintHelp()
			os.Exit(0)
		}
	}

	for flagIdx := 0; flagIdx < len(args); {
		var flagName, flagValue string
		if flagIdx+1 < len(args) {
			flagName, flagValue = args[flagIdx], args[flagIdx+1]
			flagIdx = flagIdx + 2
		} else {
			flagName = args[flagIdx]
			flagIdx = flagIdx + 1
		}

		switch strings.TrimPrefix(flagName, "--") {
		case "port":
			Port, err = strconv.Atoi(flagValue)
			if err != nil {
				return fmt.Errorf("error while parsing the port: %w", err)
			} else if Port < 1024 || Port > 65535 {
				return fmt.Errorf("incorrect range port, port must me between 1024 and 65535")
			}
		case "endpoints":
			PrintEndPoints()
			os.Exit(0)

		default:
			return fmt.Errorf("unknown flag: %s", flagName)
		}
	}

	return nil
}

func PrintHelp() {
	fmt.Println(`Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --endpoints  Show the api endpoints.
  `)
}

func PrintEndPoints() {
	fmt.Println(`
==========================================================

▶ Orders
  ├─ POST    /orders
  │          → Create a new order.
  ├─ GET     /orders
  │          → Retrieve all orders.
  ├─ GET     /orders/open
  │          → Get a list of open orders.
  ├─ GET     /orders/{id}
  │          → Retrieve a specific order by ID.
  ├─ PUT     /orders/{id}
  │          → Update an existing order.
  ├─ DELETE  /orders/{id}
  │          → Delete an order.
  ├─ POST    /orders/{id}/close
  │          → Close an order.
  └─ GET     /orders/numberOfOrderedItems
             ?startDate={startDate}&endDate={endDate}
  │          → Returns a list of ordered items and their quantities for a specified time period.
  │
  │          Parameters:
  │            - startDate (optional): Start date in YYYY-MM-DD format.
  │            - endDate   (optional): End date in YYYY-MM-DD format.

▶ Menu Items
  ├─ POST    /menu
  │          → Add a new menu item.
  ├─ GET     /menu
  │          → Retrieve all menu items.
  ├─ GET     /menu/{id}
  │          → Retrieve a specific menu item.
  ├─ PUT     /menu/{id}
  │          → Update a menu item.
  └─ DELETE  /menu/{id}
             → Delete a menu item.

▶ Inventory
  ├─ POST    /inventory
  │          → Add a new inventory item.
  ├─ GET     /inventory
  │          → Retrieve all inventory items.
  ├─ GET     /inventory/{id}
  │          → Retrieve a specific inventory item.
  ├─ PUT     /inventory/{id}
  │          → Update an inventory item.
  ├─ DELETE  /inventory/{id}
  │          → Delete an inventory item.
  └─ GET     /inventory/getLeftOvers
             ?sortBy={value}&page={page}&pageSize={pageSize}
  │          → Returns the inventory leftovers in the coffee shop, including sorting and pagination options.
  │
  │          Parameters:
  │            - sortBy   (optional): Sort method, e.g., "price" or "quantity".
  │            - page     (optional): Page number, starting from 1.
  │            - pageSize (optional): Number of items per page (default: 10).

▶ Aggregations
  ├─ GET     /reports/total-sales
  │          → Get the total sales amount.
  ├─ GET     /reports/popular-items
  │          → Get a list of popular menu items.
  ├─ GET     /reports/search
  │          ?q={query}&filter={orders|menu|all}&minPrice={minPrice}&maxPrice={maxPrice}
  │          → Search through orders, menu items, and customers with partial matching and ranking.
  │
  │          Parameters:
  │            - q         (required): Search query string.
  │            - filter    (optional): "orders", "menu", or "all" (default).
  │            - minPrice  (optional): Minimum price filter.
  │            - maxPrice  (optional): Maximum price filter.
  └─ GET     /reports/orderedItemsByPeriod
             ?period={day|month}&month={month}&year={year}
  │          → Returns the number of orders for the specified period.
  │
  │          Parameters:
  │            - period (required): "day" (group by day) or "month" (group by month).
  │            - month  (optional): Month name (e.g., "October"). Required if period=day.
  │            - year   (optional): Year. Required if period=month.

==========================================`)
}
