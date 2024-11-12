package jsonrepository

import (
	"hot-coffee/internal/001Domain/entities"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/utils"
	"os"
	"path/filepath"
)

var (
	inventoryPath = filepath.Join(flag.StoragePath, "inventory.json")
)

func initInventoryJSONRepository() {
	_, err := os.Open(inventoryPath)
	if !os.IsNotExist(err) {
		utils.FatalError("Error while opening inventory JSON file", err)
	} else if os.IsNotExist(err) {
		_, err := os.OpenFile(inventoryPath, os.O_CREATE, 0755)
		utils.FatalError("Error while creating inventory JSON file", err)
		if err == nil {
			return
		}
	}

	// Initialize repository intstance
	NewInventoryRepository()

}

// Singleton pattern
var inventoryRepositoryInstance *inventoryRepository

type inventoryRepository struct {
	storage map[string]*entities.InventoryItem
}

func NewInventoryRepository() (*inventoryRepository, error) {
	if inventoryRepositoryInstance != nil {
		return inventoryRepositoryInstance, nil
	} else {
		inventoryRepositoryInstance = &inventoryRepository{
			storage: make(map[string]*entities.InventoryItem),
		}

		return inventoryRepositoryInstance, nil
	}
}

func Create(item entities.InventoryItem) error {
	return nil
}
func GetAll() ([]entities.InventoryItem, error) {
	return nil, nil
}

func GetById(id string) (entities.InventoryItem, error) {
	return nil, nil
}
func Update(id string, item entities.InventoryItem) error {
	return nil
}
func Delete(id string) error {
	return nil
}
