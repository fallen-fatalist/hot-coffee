package jsonrepository

import (
	"encoding/json"
	"errors"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/utils"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

// Errors
var (
	ErrItemDoesntExist       = errors.New("inventory item doesn't exist by provided id")
	ErrIngridientIDDifferent = errors.New("incorrect ingridient id provided")
	ErrItemAlreadyExists     = errors.New("item already exists in inventory")
)

var (
	inventoryPath = filepath.Join(flag.StoragePath, "inventory.json")
)

// Singleton pattern
var inventoryRepositoryInstance *inventoryRepository

type inventoryRepository struct {
	storage map[string]*entities.InventoryItem
}

func NewInventoryRepository() *inventoryRepository {
	if inventoryRepositoryInstance != nil {
		return inventoryRepositoryInstance
	} else {
		inventoryRepositoryInstance = &inventoryRepository{
			storage: make(map[string]*entities.InventoryItem),
		}

		// Open file:
		inventoryPayload, err := os.ReadFile(inventoryPath)

		// File validation
		if !os.IsNotExist(err) {
			utils.FatalError("Error while opening inventory JSON file", err)
			// File does not exist
		} else if os.IsNotExist(err) {
			_, err := os.OpenFile(inventoryPath, os.O_CREATE, 0755)
			utils.FatalError("Error while creating inventory JSON file", err)
			if err == nil {
				slog.Debug("Created empty inventory JSON file")
				return inventoryRepositoryInstance
			}
		}

		inventoryRepositoryInstance.loadFromJSON(inventoryPayload)

		return inventoryRepositoryInstance
	}
}

func (r *inventoryRepository) loadFromJSON(payload []byte) error {
	// Load from file to RAM
	var items []entities.InventoryItem
	err := json.Unmarshal([]byte(payload), &items)
	if err != nil {
		log.Fatalf("Error unmarshalling Inventory JSON file: %v", err)
		os.Exit(1)
	}

	for _, item := range items {
		inventoryRepositoryInstance.storage[item.IngredientID] = &item
	}

	items = nil

	return nil
}

func (r *inventoryRepository) saveToJSON() error {
	items := make([]*entities.InventoryItem, 0, len(r.storage))
	for _, item := range r.storage {
		items = append(items, item)
	}

	// Write to JSON file
	jsonPayload, err := json.MarshalIndent(items, "", "   ")
	if err != nil {
		slog.Error("Error while Marshalling inventory items: %s", err)
		return err
	}
	items = nil
	err = os.WriteFile(inventoryPath, []byte(jsonPayload), 0755)
	if err != nil {
		slog.Error("Error while writing into %s file: %s", inventoryPath, err)
		return err
	}

	slog.Info("Inventory repository synced data with JSON file")
	return nil
}
func (r *inventoryRepository) Create(item entities.InventoryItem) error {
	if _, exists := r.storage[item.IngredientID]; exists {
		return ErrItemAlreadyExists
	}

	r.storage[item.IngredientID] = &item
	return r.saveToJSON()
}

func (r *inventoryRepository) GetAll() ([]entities.InventoryItem, error) {
	items := make([]entities.InventoryItem, 0, len(r.storage))
	for _, item := range r.storage {
		items = append(items, *item)
	}

	return items, nil
}

func (r *inventoryRepository) GetById(id string) (entities.InventoryItem, error) {
	if item, exists := r.storage[id]; exists {
		return *item, nil
	}
	return entities.InventoryItem{}, ErrItemDoesntExist
}
func (r *inventoryRepository) Update(id string, item entities.InventoryItem) error {
	if _, exists := r.storage[id]; exists {
		r.storage[id] = &item
		// Sync with file
		return r.saveToJSON()
	}

	return ErrItemDoesntExist
}
func (r *inventoryRepository) Delete(id string) error {
	if _, exists := r.storage[id]; exists {
		delete(r.storage, id)
		return r.saveToJSON()
	}
	return ErrItemDoesntExist
}
