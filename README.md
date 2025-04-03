# Hot Coffee - Coffee Shop Management System

## Overview

**Hot Coffee** is a backend system built with **Go** to manage a coffee shop's orders, menu items, and inventory. The application provides a **RESTful API** for handling key operations, with data stored in a **PostgreSQL database inside a Docker container**. It follows a **layered architecture** (Presentation, Business Logic, Data Access) for clean, maintainable code.

## Key Features

- **Order Management**: Create, update, delete, and close customer orders.
- **Menu Management**: Add, update, retrieve, and delete menu items.
- **Inventory Management**: Track and update ingredient stock levels.
- **Reports**: View different aggregations.
- **Logging**: Logs events and errors for monitoring and debugging.

## Architecture

The system uses a **three-layer architecture**:
- **Core**: Contains the entities manipulated in Services and Repositories layer
- **Presentation**: Manage HTTP requests and responses.
- **Services**: Contain core business logic.
- **Repositories**: Handle data storage and retrieval from JSON files.

## Dependencies
- **Git**
- **Golang 1.23.2 Compiler**
- **Docker Compose**
- **Make**

## Running the Application

1. Clone the repository:
   ```bash
   git clone https://github.com/fallen-fatalist/hot-coffee.git
   cd hot-coffee
   ```

2. Up the application:  
   ```bash 
   make up    # Ups docker containers with application and database 
   ```
3. If needed, you can stop the application or check logs using:  
   ```bash
   make down  # Stop containers  
   make logs  # View application logs  
   ```

## Program variables
* The application will start a server on the default port (or use `--port` to specify a different one).
* To get help:  
```bash
go run main.go --help
```
![image](https://github.com/user-attachments/assets/e6a36a4f-6ec5-454e-958c-fc019a95e109)

* To get list of endpoints:
```bash
go run main.go --endpoints
```
![image](https://github.com/user-attachments/assets/f47e5e38-9a3f-414f-bc5a-92d3051e4a21)

## Project Structure

```
.
├── cmd                                          # Project initializer
│   ├── app.go
│   └── routes.go
├── db_init                                      # Postgres migration scripts
│   ├── 001_create_customers.sql
│   ├── 002_create_orders.sql
│   ├── 003_create_order_status_history.sql
│   ├── 004_create_menu_items.sql
│   ├── 005_create_order_items.sql
│   ├── 006_create_price_history.sql
│   ├── 007_create_inventory.sql
│   ├── 008_create_menu_items_ingredients.sql
│   ├── 010_create_inventory_transactions.sql
│   ├── 011_mock_customers.sql
│   ├── 012_mock_orders.sql
│   ├── 013_mock_menu_items.sql
│   ├── 014_mock_inventory.sql
│   ├── 015_mock_menu_item_ingredients.sql
│   ├── 016_mock_inventory_transactions.sql
│   ├── 017_mock_order_items.sql
│   ├── 018_mock_order_status_history.sql
│   ├── 019_mock_price_history.sql
│   ├── 020_mock_inventory_transactions.sql
│   ├── 021_create_index_orders.sql
│   ├── 022_create_index_order_items.sql
│   └── 023_create_index_menu_items_ingredients.sql
├── docker-compose.yml
├── Dockerfile
├── docs                                        # Project documentation
│   ├── api
│   │   └── frapuccino.postman_collection.json
│   └── erd
│       ├── Database ERD Frapuccino.pdf
│       └── erd.jpg
├── go.mod
├── go.sum
├── internal
│   ├── core                                    # Core layer
│   │   ├── entities
│   │   │   ├── aggregation.go
│   │   │   ├── inventory_item.go
│   │   │   ├── menu_item.go
│   │   │   └── order.go
│   │   └── errors
│   │       └── errors.go
│   ├── dto
│   │   └── order_dto.go
│   ├── flag
│   │   └── flag.go
│   ├── infrastructure                          # Infrastructure which implements services  
│   │   ├── server                              # Presentation layer
│   │   │   └── http
│   │   │       ├── aggregation_handler.go
│   │   │       ├── helpers.go
│   │   │       ├── inventory_handler.go
│   │   │       ├── menu_handler.go
│   │   │       ├── middleware.go
│   │   │       └── order_handler.go
│   │   └── storage                             # Repository implementation
│   │       └── postgres
│   │           ├── inventory_repository.go
│   │           ├── menu_repository.go
│   │           ├── order_repository.go
│   │           └── storage.go
│   ├── repository                              # Repository interfaces
│   │   └── repository.go
│   ├── service                                 # Service layer
│   │   ├── service.go
│   │   └── serviceinstance
│   │       ├── aggregation_service.go
│   │       ├── inventory_service.go
│   │       ├── menu_service.go
│   │       ├── order_service.go
│   │       ├── service.go
│   │       └── validator.go
│   ├── utils
│   │   └── utils.go
│   └── vo
│       ├── inventory_vo.go
│       └── order_vo.go
├── main.go                                    # Project entrypoint
├── Makefile                                   # Build scripts
├── README.md
└── TODO.md                                       
```

## API Endpoints

### **Orders**
- `GET /orders` - Retrieve all orders.  
- `GET /orders/open` - Get open orders.  
- `POST /orders` – Create an order.  
- `GET /orders/{id}` – Get an order.  
- `PUT /orders/{id}` – Update an order.  
- `DELETE /orders/{id}` – Delete an order.  
- `POST /orders/{id}/close` – Close an order.
- `GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}` - Number of ordered items.
- `POST /orders/batch-process` - Bulk order processing.  

### **Menu**
- `GET /menu` - Retrieve all menu items.  
- `POST /menu` – Add a menu item.  
- `GET /menu/{id}` – Get a menu item.  
- `PUT /menu/{id}` – Update a menu item.  
- `DELETE /menu/{id}` – Delete a menu item.  

### **Inventory**
- `GET /inventory` - Retrieve all inventory items.  
- `POST /inventory` – Add an inventory item.  
- `GET /inventory/{id}` – Get an inventory item.  
- `PUT /inventory/{id}` – Update an inventory item.  
- `DELETE /inventory/{id}` – Delete an inventory item.
- `GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}` - Get leftovers.

### **Reports**
- `GET /reports/total-sales` – Total sales.  
- `GET /reports/popular-items` – Popular menu items.  
- `GET /reports/search?q={searchQuery}&filter={filter}&minPrice={minPrice}&maxPrice={maxPrice}` - Full text search report.  
- `GET /reports/orderedItemsByPeriod?period={day|month}&month={month}` - Ordered items by period.  
  
 

## Request/Response Examples:

**Now examples can be seen in Postman collection under /docs/api folder**
**Further the examples will be added here**

## Data Storage

Data is stored in **PostgreSQL**, with the following tables:

- `customers` – Stores customer information.
- `orders` – Stores customer orders.
- `order_status_history` – Tracks changes in order statuses.
- `menu_items` – Stores menu items (products).
- `order_items` – Tracks items in each order.
- `price_history` – Stores the price history for menu items.
- `inventory` – Tracks ingredient stock and prices.
- `menu_items_ingredients` – Stores the relationship between menu items and their ingredients.
- `inventory_transactions` – Tracks inventory changes related to orders.

### ERD diagram
![image](https://github.com/user-attachments/assets/d2c85a88-a5c2-41f9-aaeb-2bdde292248e)

