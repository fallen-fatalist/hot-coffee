package jsonrepository

import (
	"fmt"
	"log/slog"
	"os"

	"hot-coffee/internal/flag"
)

func Init() {
	// Validate or initalize data path
	_, err := os.ReadDir(flag.StoragePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error while opening data path: %s", err)
		os.Exit(1)
	} else if os.IsNotExist(err) {
		err := os.Mkdir(flag.StoragePath, 0o755)
		if err != nil {
			fmt.Printf("Error while creating data storage: %s", err)
			os.Exit(1)
		} else {
			slog.Info(fmt.Sprintf("Created data path in: %s", flag.StoragePath))
		}
	}
	// Initialize other repositories
	NewInventoryRepository()
	NewMenuRepository()
	NewOrderRepository()
	slog.Info("All JSON repositories are initialized")
}
