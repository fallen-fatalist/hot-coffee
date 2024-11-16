# Hot Coffee - Coffee Shop Management System

## Overview

**Hot Coffee** is a backend system built with **Go** to manage a coffee shop's orders, menu items, and inventory. The application provides a **RESTful API** for handling key operations, with data stored in **JSON files** locally. It follows a **layered architecture** (Presentation, Business Logic, Data Access) for clean, maintainable code.

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

## Running the Application

1. Clone the repository:
   ```bash
   git clone https://github.com/fallen-fatalist/hot-coffee.git
   cd hot-coffee
   ```

2. Run the application:
   ```bash
   go run main.go
   ```
   or
   ```
   go build -o <binary name> .
   ./<binary name>
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
├── cmd
│   ├── app.go                          # main entry point
│   └── routes.go                       # contains router initialization
├── data                                # storage directory
│   ├── inventory.json
│   ├── menu_items.json
│   └── orders.json
├── go.mod
├── hot-coffee.md
├── internal                            # main logic
│   ├── core                            # main entities 
│   │   └── entities
│   │       ├── inventory_item.go
│   │       ├── menu_item.go
│   │       └── order.go
│   ├── flag                            # flag handling
│   │   └── flag.go
│   ├── infrastructure
│   │   └── controllers                 # controllers
│   │       ├── aggregation_handler.go
│   │       ├── inventory_handler.go
│   │       ├── menu_handler.go
│   │       └── order_handler.go
│   ├── repositories                    # repositories
│   │   ├── jsonrepository              # repositories instance implementation
│   │   │   ├── inventory_repository.go
│   │   │   ├── menu_repository.go
│   │   │   ├── order_repository.go
│   │   │   └── repository.go.go
│   │   └── repository                   # repositories interfaces
│   │       ├── inventory_repository.go
│   │       ├── menu_repository.go
│   │       └── order_repository.go
│   ├── services                         # application services (domain services)
│   │   ├── service                      # service interfaces
│   │   │   ├── inventory_service.go
│   │   │   ├── menu_service.go
│   │   │   └── order_service.go
│   │   └── serviceinstance              # services instance
│   │       ├── inventory_service.go
│   │       ├── menu_service.go
│   │       ├── order_service.go
│   │       └── service.go
│   └── utils                            # utilities reused in the project 
│       └── utils.go
├── main.go
├── README.md
```

## API Endpoints

- **Orders**:
  - `GET /orders` - Retrieve all orders.
  - `GET /orders/open` - Get open orders
  - `POST /orders` – Create an order.
  - `GET /orders/{id}` – Get an order.
  - `PUT /orders/{id}` – Update an order.
  - `DELETE /orders/{id}` – Delete an order.
  - `POST /orders/{id}/close` – Close an order.

- **Menu**:
  - `GET /menu` - Retrieve all menu items.
  - `POST /menu` – Add a menu item.
  - `GET /menu/{id}` – Get a menu item.
  - `PUT /menu/{id}` – Update a menu item.
  - `DELETE /menu/{id}` – Delete a menu item.

- **Inventory**:
  - `GET /inventory` - Retrieve all inventory items
  - `POST /inventory` – Add an inventory item. 
  - `GET /inventory/{id}` – Get an inventory item.
  - `PUT /inventory/{id}` – Update an inventory item.
  - `DELETE /inventory/{id}` – Delete an inventory item.

- **Reports**:
  - `GET /reports/total-sales` – Total sales.
  - `GET /reports/popular-items` – Popular menu items.
 

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

Data is stored in **JSON** files (by default in the `data` folder):

- `orders.json` – Stores customer orders.
- `menu_items.json` – Stores menu items (product, ingredients).
- `inventory.json` – Tracks ingredient stock.





## Error Handling

- **400 Bad Request** for invalid input.
- **404 Not Found** when resources are not found.
- **500 Internal Server Error** for unexpected server issues.

## Logging
The application uses Go's `log/slog` package to log significant events and errors.
