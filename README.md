# Hot Coffee - Coffee Shop Management System

## Overview

**Hot Coffee** is a backend system built with **Go** to manage a coffee shop's orders, menu items, and inventory. The application provides a **RESTful API** for handling key operations, with data stored in a **PostgreSQL database inside a Docker container**. It follows a **layered architecture** (Presentation, Business Logic, Data Access) for clean, maintainable code.

## Key Features

- **Order Management**: Create, update, delete, and close customer orders.
- **Menu Management**: Add, update, retrieve, and delete menu items.
- **Inventory Management**: Track and update ingredient stock levels.
- **Sales Reports**: View total sales and popular menu items.
- **Logging**: Logs events and errors for monitoring and debugging.

## Architecture

The system uses a **three-layer architecture**:
- **Core**: Contains the entities manipulated in Services and Repositories layer
- **Controllers**: Manage HTTP requests and responses.
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

2. Run the application:  
   ```bash
   make build  
   make up  
   ```
If needed, you can stop the application or check logs using:  
   ```bash
   make down  # Stop and remove containers  
   make logs  # View application logs  
   ```

## Program variables
* The application will start a server on the default port (or use `--port` to specify a different one).
* To get help:  
```
go run main.go --help
```
![image](https://github.com/user-attachments/assets/72ed27db-1afa-4044-862d-ae867eed22d5)

* To get list of endpoints:
```
go run main.go --endpoints
```
![Uploading image.png…]()

* To change directory where save data:
```
go run main.go --dir <my-dir>
```

## Project Structure

```
hot-coffee/
├── cmd/                                # Application entry points
│   ├── app.go                          # Main entry point
│   ├── routes.go                        # Route definitions
├── db_init/                            # Database initialization scripts
├── erd/                                # Entity-Relationship Diagrams
├── internal/                           # Core application logic
│   ├── core/                           # Core domain logic
│   │   ├── entities/                   # Domain entities
│   │   │   ├── aggregation.go
│   │   │   ├── inventory_item.go
│   │   │   ├── menu_item.go
│   │   │   ├── order.go
│   │   ├── errors/                     # Error handling
│   │   │   ├── errors.go
│   │   ├── dto/                        # Data Transfer Objects
│   │   │   ├── order_dto.go
│   ├── flag/                           # Command-line flag handling
│   │   ├── flag.go
│   ├── infrastructure/                 # Infrastructure components
│   │   ├── server/                     # Server setup (currently empty)
│   │   ├── storage/                    # Storage implementations
│   │   │   ├── postgres/               # PostgreSQL storage implementation
│   │   │   │   ├── inventory_repository.go
│   │   │   │   ├── menu_repository.go
│   │   │   │   ├── order_repository.go
│   │   │   │   ├── storage.go
│   ├── repository/                     # Repository interfaces
│   │   ├── repository.go
│   ├── service/                        # Business logic layer
│   │   ├── serviceinstance/            # Service implementations
│   │   │   ├── aggregation_service.go
│   │   │   ├── inventory_service.go
│   │   │   ├── menu_service.go
│   │   │   ├── order_service.go
│   │   │   ├── service.go
│   │   │   ├── validator.go
│   │   ├── service.go
│   ├── utils/                          # Utility functions
│   │   ├── utils.go
│   ├── vo/                             # Value Objects
│   │   ├── inventory_vo.go
│   │   ├── order_vo.go
├── .gitignore                          # Git ignore file
├── docker-compose.yml                   # Docker Compose configuration
├── Dockerfile                          # Docker configuration
├── go.mod                              # Go module definition
├── go.sum                              # Go module checksums
├── main.go                             # Main application entry point
├── Makefile                            # Build automation
├── README.md                           # Project documentation
├── TODO.md                             # Task list and future enhancements
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

### **Reports**
- `GET /reports/total-sales` – Total sales.  
- `GET /reports/popular-items` – Popular menu items.  

### **New Endpoints**
- `GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}` - Number of ordered items.  
- `GET /reports/search?q={searchQuery}&filter={filter}&minPrice={minPrice}&maxPrice={maxPrice}` - Full text search report.  
- `GET /reports/orderedItemsByPeriod?period={day|month}&month={month}` - Ordered items by period.  
- `GET /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}` - Get leftovers.  
- `POST /orders/batch-process` - Bulk order processing.  

## Request/Response Examples:
* Create Order Request:
```
POST /orders
Content-Type: application/json

{
  "customer_name": "John Doe",
  "items": [
    {
      "product_id": "espresso",
      "quantity": 2
    },
    {
      "product_id": "croissant",
      "quantity": 1
    }
  ]
}
```
* Create menu request
```
POST /menu
Content-Type: application/json
 {
      "product_id": "muffin",
      "name": "Blueberry Muffin",
      "description": "Freshly baked muffin with blueberries",
      "price": 2,
      "ingredients": [
         {
            "ingredient_id": "flour",
            "quantity": 100
         },
         {
            "ingredient_id": "blueberries",
            "quantity": 20
         },
         {
            "ingredient_id": "sugar",
            "quantity": 30
         }
      ]
}
```
* Create inventory item request:
```
POST /inventory
Content-Type: application/json
{
      "ingredient_id": "espresso_shot",
      "name": "Espresso Shot",
      "quantity": 500,
      "unit": "shots"
}
```
* Total Sales Aggregation Response:
```
HTTP/1.1 200 OK
Content-Type: application/json

{
  "total_sales": 1500.50
}
```
* Popular Items Aggregation Response:
```
HTTP/1.1 200 OK
Content-Type: application/json
[
   {
      "product_id": "espresso",
      "product_name": "Espresso",
      "total_sales_count": 3
   }
]
```

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





## Error Handling

- **400 Bad Request** for invalid input.
- **404 Not Found** when resources are not found.
- **500 Internal Server Error** for unexpected server issues.

## Logging
The application uses Go's `log/slog` package to log significant events and errors.
