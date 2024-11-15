package jsonrepository

import (
	"encoding/json"
	"errors"
	"fmt"
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
	ErrInventoryItemDoesntExist   = errors.New("inventory item doesn't exist by provided id")
	ErrIngridientIDDifferent      = errors.New("incorrect ingridient id provided")
	ErrInventoryItemAlreadyExists = errors.New("item already exists in inventory")
)

// Singleton pattern
var inventoryRepositoryInstance *inventoryRepository

type inventoryRepository struct {
	repository         map[string]*entities.InventoryItem
	repositoryFileName string
}

func NewInventoryRepository() *inventoryRepository {
	if inventoryRepositoryInstance != nil {
		return inventoryRepositoryInstance
	} else {
		inventoryRepositoryInstance = &inventoryRepository{
			repository:         make(map[string]*entities.InventoryItem),
			repositoryFileName: filepath.Join(flag.StoragePath, "inventory.json"),
		}

		// Open file:
		inventoryPayload, err := os.ReadFile(inventoryRepositoryInstance.repositoryFileName)

		// File validation
		if !os.IsNotExist(err) {
			utils.FatalError("Error while opening inventory JSON file", err)
			// File does not exist
		} else if os.IsNotExist(err) {
			_, err := os.OpenFile(inventoryRepositoryInstance.repositoryFileName, os.O_CREATE, 0o755)
			fillJSONWithArray(inventoryRepositoryInstance.repositoryFileName)
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
		inventoryRepositoryInstance.repository[item.IngredientID] = &item
	}

	items = nil

	return nil
}

func fillJSONWithArray(fileName string) {
	err := os.WriteFile(fileName, []byte("[]"), 0o755)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while filling %s JSON file: %s", fileName, err))
		os.Exit(1)
	}
}

func (r *inventoryRepository) saveToJSON() error {
	items := make([]*entities.InventoryItem, 0, len(r.repository))
	for _, item := range r.repository {
		items = append(items, item)
	}

	// Write to JSON file
	jsonPayload, err := json.MarshalIndent(items, "", "   ")
	if err != nil {
		slog.Error(fmt.Sprintf("Error while Marshalling inventory items: %s", err))
		return err
	}
	items = nil
	err = os.WriteFile(r.repositoryFileName, []byte(jsonPayload), 0755)
	if err != nil {
		slog.Error("Error while writing into %s file: %s", r.repositoryFileName, err)
		return err
	}

	slog.Info("Inventory repository synced data with JSON file")
	return nil
}
func (r *inventoryRepository) Create(item entities.InventoryItem) error {
	if _, exists := r.repository[item.IngredientID]; exists {
		return ErrInventoryItemAlreadyExists
	}

	r.repository[item.IngredientID] = &item
	return r.saveToJSON()
}

func (r *inventoryRepository) GetAll() ([]entities.InventoryItem, error) {
	items := make([]entities.InventoryItem, 0, len(r.repository))
	for _, item := range r.repository {
		items = append(items, *item)
	}

	return items, nil
}

func (r *inventoryRepository) GetById(id string) (entities.InventoryItem, error) {
	if item, exists := r.repository[id]; exists {
		return *item, nil
	}
	return entities.InventoryItem{}, ErrInventoryItemDoesntExist
}
func (r *inventoryRepository) Update(id string, item entities.InventoryItem) error {
	if _, exists := r.repository[id]; exists {
		r.repository[id] = &item
		// Sync with file
		return r.saveToJSON()
	}

	return ErrInventoryItemDoesntExist
}
func (r *inventoryRepository) Delete(id string) error {
	if _, exists := r.repository[id]; exists {
		delete(r.repository, id)
		return r.saveToJSON()
	}
	return ErrInventoryItemDoesntExist
}
