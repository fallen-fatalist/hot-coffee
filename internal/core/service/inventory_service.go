package service

import "net/http"

// service variable
var (
	InventoryService = &inventoryService{}
)

type inventoryService struct {
	// Repository insertion
	//inventoryRepository
}

func (s *inventoryService) PostInventoryItem(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *inventoryService) GetInventory(w http.ResponseWriter, r *http.Request) {
	return
}
