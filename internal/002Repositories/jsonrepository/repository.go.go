package jsonrepository

import (
	"fmt"
	"hot-coffee/internal/flag"
	"log/slog"
	"os"
)

func Init() {
	// Validate or initalize data path
	_, err := os.ReadDir(flag.StoragePath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error while opening data path: %s", err)
		os.Exit(1)
	} else if os.IsNotExist(err) {
		err := os.Mkdir(flag.StoragePath, 0755)
		if err != nil {
			fmt.Printf("Error while creating data storage: %s", err)
			os.Exit(1)
		} else {
			slog.Info(fmt.Sprintf("Created data path in: %s", flag.StoragePath))
		}
	}
	// Initialize other repositories
	initInventoryJSONRepository()
	initMenuJSONRepository()
	initOrderJSONRepository()
	slog.Info("All JSON repositories are initialized")
}