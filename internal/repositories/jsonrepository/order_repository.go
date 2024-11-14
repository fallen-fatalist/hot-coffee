package jsonrepository

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/utils"
)

type orderRepository struct {
	repository         map[string]*entities.Order
	repositoryFilename string
}

// Singleton pattern
var orderRepositoryInstance *orderRepository

func NewOrderRepository() *orderRepository {
	if orderRepositoryInstance != nil {
		return orderRepositoryInstance
	}
	orderRepositoryInstance = &orderRepository{
		repository:         make(map[string]*entities.Order),
		repositoryFilename: filepath.Join(flag.StoragePath, "orders.json"),
	}

	// Open file:
	orderPayload, err := os.ReadFile(orderRepositoryInstance.repositoryFilename)

	// File validation
	if !os.IsNotExist(err) {
		utils.FatalError("Error while opening order JSON file", err)
		// File does not exist
	} else if os.IsNotExist(err) {
		_, err := os.OpenFile(orderRepositoryInstance.repositoryFilename, os.O_CREATE, 0o755)
		utils.FatalError("Error while creating order JSON file", err)
		if err == nil {
			slog.Debug("Created empty menu JSON file")
			return orderRepositoryInstance
		}
	}

	orderRepositoryInstance.loadFromJSON(orderPayload)

	return orderRepositoryInstance
}

func (o *orderRepository) loadFromJSON(payload []byte) error {
	// Load from file to RAM
	var orders []entities.Order
	err := json.Unmarshal([]byte(payload), &orders)
	if err != nil {
		log.Fatalf("Error unmarshalling menu JSON file: %v", err)
		os.Exit(1)
	}

	for _, order := range orders {
		orderRepositoryInstance.repository[order.ID] = &order
	}
	orders = nil
	return nil
}

func (m *orderRepository) saveToJSON() error {
	orders := make([]*entities.Order, 0, len(m.repository))
	for _, order := range m.repository {
		orders = append(orders, order)
	}

	// Write to JSON file
	jsonPayload, err := json.MarshalIndent(orders, "", "   ")
	if err != nil {
		slog.Error(fmt.Sprintf("Error while Marshalling orders: %s", err))
		return err
	}
	orders = nil
	err = os.WriteFile(m.repositoryFilename, []byte(jsonPayload), 0o755)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while writing into %s file: %s", inventoryFilename, err))
		return err
	}

	slog.Info("Orders repository synced data with JSON file")
	return nil
}

func (o *orderRepository) Create(order entities.Order) error {
	if _, exists := o.repository[order.ID]; exists {
		return nil
	}
	o.repository[order.ID] = &order
	return o.saveToJSON()
}

func (o *orderRepository) GetAll() ([]entities.Order, error) {
	orders := make([]entities.Order, 0, len(o.repository))
	for _, order := range o.repository {
		orders = append(orders, *order)
	}
	return orders, nil
}

func (o *orderRepository) GetById(id string) (entities.Order, error) {
	if _, exists := o.repository[id]; exists {
		return *o.repository[id], nil
	}
	return entities.Order{}, ErrMenuItemDoesntExist
}

func (o *orderRepository) Update(id string, order entities.Order) error {
	if id != order.ID {
		
	}
	if _, exists := o.repository[id]; exists {
		o.repository[id] = &order
		return o.saveToJSON()
	}
	return ErrMenuItemDoesntExist
}

func (o *orderRepository) Delete(id string) error {
	if _, exists := o.repository[id]; exists {
		delete(o.repository, id)
		return o.saveToJSON()
	}
	return ErrMenuItemDoesntExist
}
