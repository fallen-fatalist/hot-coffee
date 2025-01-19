package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"hot-coffee/internal/repository"
	"log/slog"
	"time"
)

// Database credentials
const (
	dbName   = "frappuccino"
	host     = "db"
	userName = "latte"
	password = "latte"
	port     = "5432"
)

// postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName
var dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", userName, password, host, port, dbName)

// Singleton Pattern: SQL database connection
var PostgresDB *sql.DB

func NewRepository() *repository.Repository {
	return &repository.Repository{
		Inventory: NewInventoryRepository(),
		Menu:      nil,
		Order:     nil,
	}
}

func openDB() (*sql.DB, error) {
	if PostgresDB != nil {
		return PostgresDB, nil
	}
	slog.Info("Trying to connect to PostgreSQL database...")
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	PostgresDB = db
	slog.Info("PostgreSQL database connection established")
	// Return the sql.DB connection pool.
	return PostgresDB, nil
}
