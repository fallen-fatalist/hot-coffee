package service

import (
	"fmt"
	"hot-coffee/internal/002Repositories/repositories"
	"net/http"
)

// service variable
var (
	InventoryService = &inventoryService{}
)

// errors
var (
	ErrNilRepository = "Nil repository provided into constructor"
)

type inventoryService struct {
	InventoryRepository repositories.InventoryRepository
}

func NewInventoryService(repository repositories.InventoryRepository) (*inventoryService, error) {
	if repository == nil {
		return nil, fmt.Errorf("Error while creating Inventory service: %w", ErrNilRepository)
	}
	return &inventoryService{repository}, nil
}

func (s *inventoryService) PostInventoryItem(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *inventoryService) GetInventory(w http.ResponseWriter, r *http.Request) {
	return
}
