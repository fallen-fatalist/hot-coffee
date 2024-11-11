package handler

import (
	"hot-coffee/internal/core/service"
	"net/http"
)

func HandleInventory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		service.InventoryService.GetInventory(w, r)
	case http.MethodPost:
		service.InventoryService.PostInventoryItem(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func HandleInventoryItem(w http.ResponseWriter, r *http.Request) {
	// ID validation
	switch r.Method {
	case http.MethodGet:
	case http.MethodPut:
	case http.MethodDelete:
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
