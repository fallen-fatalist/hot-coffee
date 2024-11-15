package jsonrepository

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/flag"
	"hot-coffee/internal/utils"
)

// Errors
var (
	ErrOrderAlreadyExists = errors.New("order with such id already exists")
	ErrOrderIsNotExist    = errors.New("order with provided id does not exist")
	ErrUUIDGeneration     = errors.New("error while uuid generation")
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
	order.CreatedAt = time.Now().String()
	var err error
	for {
		order.ID, err = generateOrderID()
		if err != nil {
			return err
		}

		if _, exists := o.repository[order.ID]; !exists {
			break
		}
	}
	o.repository[order.ID] = &order
	return o.saveToJSON()
}

func generateOrderID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", ErrUUIDGeneration
	}

	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return "order-" + uuid, nil
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
	return entities.Order{}, ErrOrderIsNotExist
}

func (o *orderRepository) Update(id string, order entities.Order) error {
	if _, exists := o.repository[id]; exists {
		o.repository[id].CustomerName = order.CustomerName
		o.repository[id].Items = order.Items
		o.repository[id].Status = order.Status
		return o.saveToJSON()
	}
	return ErrOrderIsNotExist
}

func (o *orderRepository) Delete(id string) error {
	if _, exists := o.repository[id]; exists {
		delete(o.repository, id)
		return o.saveToJSON()
	}
	return ErrOrderIsNotExist
}
