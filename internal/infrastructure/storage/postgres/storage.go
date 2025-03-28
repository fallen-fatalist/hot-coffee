package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"hot-coffee/internal/repository"
	"log/slog"
	"os"
	"time"
)

// Database credentials
var (
	dbName   = os.Getenv("DB_NAME")
	host     = os.Getenv("DB_HOST")
	userName = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	port     = os.Getenv("DB_PORT")
)

// postgres://YourUserName:YourPassword@YourHostname:5432/YourDatabaseName
// postgres://username:password@host/dbname
var dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", userName, password, host, port, dbName)

// Singleton Pattern: SQL database connection
var postgresDB *sql.DB

func NewRepository() *repository.Repository {
	return &repository.Repository{
		Inventory: NewInventoryRepository(),
		Menu:      NewMenuRepository(),
		Order:     NewOrderRepository(),
	}
}

func openDB() (*sql.DB, error) {

	if postgresDB != nil {
		return postgresDB, nil
	}

	slog.Info("Trying to connect to PostgreSQL database...")

	// Retry logic: attempt to connect multiple times
	maxRetries := 6                  // Try 6 times (30 seconds total if we wait 5 seconds between retries)
	retryInterval := 5 * time.Second // Retry every 5 seconds
	var db *sql.DB
	var err error

	// Try connecting up to maxRetries times
	for range maxRetries {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			// Create a context with a timeout for the Ping operation
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Try to ping the database to check if it's available
			err = db.PingContext(ctx)
			if err == nil {
				// If the ping is successful, use the connection
				postgresDB = db
				slog.Info("PostgreSQL database connection established")
				return postgresDB, nil
			}
		}

		// If any error occurs, log it and retry after a delay
		//slog.Errorf("Failed to connect to PostgreSQL, retrying in %v... (attempt %d/%d)", retryInterval, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	slog.Error("Failed to connect to PostgreSQL after", "attempts", maxRetries, "interval", retryInterval, "error", err)
	return nil, err
}
