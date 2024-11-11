package service

import "net/http"

// service variable
var (
	MenuService = &menuService{}
)

type menuService struct {
	// Repository insertion
	//inventoryRepository
}

func (s *menuService) PostMenu(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *menuService) GetMenu(w http.ResponseWriter, r *http.Request) {
	return
}
