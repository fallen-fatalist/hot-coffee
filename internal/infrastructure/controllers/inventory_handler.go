package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/services/serviceinstance"
	"hot-coffee/internal/utils"
)

// Route: /inventory
func HandleInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.InventoryService.GetInventoryItems()
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}

		jsonPayload, err := json.MarshalIndent(items, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPost:
		var item entities.InventoryItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		err = serviceinstance.InventoryService.CreateInventoryItem(item)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /inventory/{id}
func HandleInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		item, err := serviceinstance.InventoryService.GetInventoryItem(id)
		if err != nil {
			utils.JSONErrorRespond(w, err)
		}

		jsonPayload, err := json.MarshalIndent(item, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPut:
		var item entities.InventoryItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		err = serviceinstance.InventoryService.UpdateInventoryItem(id, item)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		return
	case http.MethodDelete:
		err := serviceinstance.InventoryService.DeleteInventoryItem(id)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
