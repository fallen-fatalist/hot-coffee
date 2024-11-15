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
- **Handlers**: Manage HTTP requests and responses.
- **Services**: Contain core business logic.
- **Repositories**: Handle data storage and retrieval from JSON files.

### Project Structure

```
hot-coffee/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── handler/             # HTTP request handlers
│   ├── service/             # Business logic
│   └── dal/                 # Data storage logic
├── models/                  # Data models
├── data/                    # JSON files (orders, menu items, inventory)
├── go.mod                   # Go module
└── README.md                # Project documentation
```

## API Endpoints

- **Orders**: 
  - `POST /orders` – Create an order.
  - `GET /orders/{id}` – Get an order.
  - `PUT /orders/{id}` – Update an order.
  - `DELETE /orders/{id}` – Delete an order.
  - `POST /orders/{id}/close` – Close an order.

- **Menu**: 
  - `POST /menu` – Add a menu item.
  - `GET /menu/{id}` – Get a menu item.
  - `PUT /menu/{id}` – Update a menu item.
  - `DELETE /menu/{id}` – Delete a menu item.

- **Inventory**: 
  - `POST /inventory` – Add an inventory item.
  - `GET /inventory/{id}` – Get an inventory item.
  - `PUT /inventory/{id}` – Update an inventory item.
  - `DELETE /inventory/{id}` – Delete an inventory item.

- **Reports**:
  - `GET /reports/total-sales` – Total sales.
  - `GET /reports/popular-items` – Popular menu items.
  - `GET /reports/open` - Get open orders

## Data Storage

Data is stored in **JSON** files in the `data/` folder:

- `orders.json` – Stores customer orders.
- `menu_items.json` – Stores menu items (product, ingredients).
- `inventory.json` – Tracks ingredient stock.

## Requirements

- **Go 1.18+**
- **JSON Files** for data storage (no database).

## Running the Application

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/hot-coffee.git
   cd hot-coffee
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/main.go
   ```

The application will start a server on the default port (or use `--port` to specify a different one).

## Error Handling

- **400 Bad Request** for invalid input.
- **404 Not Found** when resources are not found.
- **500 Internal Server Error** for unexpected issues.

## Logging

The application uses Go's `log/slog` package to log significant events and errors.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

