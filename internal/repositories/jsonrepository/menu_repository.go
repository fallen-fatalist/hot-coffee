package jsonrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/utils"
)

// Errors
var (
	ErrMenuItemAlreadyExists = errors.New("menu item already exists")
	ErrMenuItemDoesntExist   = errors.New("menu item does not exist")
)

type menuRepository struct {
	repository         map[string]*entities.MenuItem
	repositoryFilename string
}

// Singleton pattern
var menuRepositoryInstance *menuRepository

func NewMenuRepository() *menuRepository {
	if menuRepositoryInstance != nil {
		return menuRepositoryInstance
	} else {
		menuRepositoryInstance = &menuRepository{
			repository:         make(map[string]*entities.MenuItem),
			repositoryFilename: filepath.Join(flag.StoragePath, "menu_items.json"),
		}

		// Open file:
		menuPayload, err := os.ReadFile(menuRepositoryInstance.repositoryFilename)

		// File validation
		if !os.IsNotExist(err) {
			utils.FatalError("Error while opening menu JSON file", err)
			// File does not exist
		} else if os.IsNotExist(err) {
			_, err := os.OpenFile(menuRepositoryInstance.repositoryFilename, os.O_CREATE, 0o755)
			utils.FatalError("Error while creating menu JSON file", err)
			if err == nil {
				slog.Debug("Created empty menu JSON file")
				return menuRepositoryInstance
			}
		}

		menuRepositoryInstance.loadFromJSON(menuPayload)

		return menuRepositoryInstance
	}
}

func (m *menuRepository) loadFromJSON(payload []byte) error {
	// Load from file to RAM
	var items []entities.MenuItem
	err := json.Unmarshal([]byte(payload), &items)
	if err != nil {
		log.Fatalf("Error unmarshalling menu JSON file: %v", err)
		os.Exit(1)
	}

	for _, item := range items {
		menuRepositoryInstance.repository[item.ID] = &item
	}
	items = nil

	return nil
}

func (m *menuRepository) saveToJSON() error {
	items := make([]*entities.MenuItem, 0, len(m.repository))
	for _, item := range m.repository {
		items = append(items, item)
	}

	// Write to JSON file
	jsonPayload, err := json.MarshalIndent(items, "", "   ")
	if err != nil {
		slog.Error(fmt.Sprintf("Error while Marshalling menu items: %s", err))
		return err
	}
	items = nil
	err = os.WriteFile(m.repositoryFilename, []byte(jsonPayload), 0o755)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while writing into %s file: %s", inventoryFilename, err))
		return err
	}

	slog.Info("Menu repository synced data with JSON file")
	return nil
}

func (m *menuRepository) Create(item entities.MenuItem) error {
	if _, exists := m.repository[item.ID]; exists {
		return ErrMenuItemAlreadyExists
	}
	m.repository[item.ID] = &item
	return m.saveToJSON()
}

func (m *menuRepository) GetAll() ([]entities.MenuItem, error) {
	items := make([]entities.MenuItem, 0, len(m.repository))
	for _, item := range m.repository {
		items = append(items, *item)
	}
	return items, nil
}

func (m *menuRepository) GetById(id string) (entities.MenuItem, error) {
	if _, exists := m.repository[id]; exists {
		return *m.repository[id], nil
	}
	return entities.MenuItem{}, ErrMenuItemDoesntExist
}

func (m *menuRepository) Update(id string, item entities.MenuItem) error {
	if _, exists := m.repository[id]; exists {
		m.repository[id] = &item
		return m.saveToJSON()
	}
	return ErrMenuItemDoesntExist
}

func (m *menuRepository) Delete(id string) error {
	if _, exists := m.repository[id]; exists {
		delete(m.repository, id)
		return m.saveToJSON()
	}
	return ErrMenuItemDoesntExist
}
