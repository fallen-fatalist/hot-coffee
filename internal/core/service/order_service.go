package service

import "net/http"

// service variable
var (
	OrderService = &orderService{}
)

type orderService struct {
	// Repository insertion
	//inventoryRepository
}

func (s *orderService) PostOrder(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *orderService) GetOrder(w http.ResponseWriter, r *http.Request) {
	return
}
