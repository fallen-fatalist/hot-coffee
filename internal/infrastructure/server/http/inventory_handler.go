package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/service/serviceinstance"
)

// Errors
var (
	ErrNonIntegerPageSize = errors.New("page size must be an integer")
	ErrNonIntegerPage     = errors.New("page must be an integer")
)

// Route: /inventory
func HandleInventory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.InventoryService.GetInventoryItems()
		if err != nil {
			if errors.Is(err, serviceinstance.ErrNoInventoryItems) {
				jsonMessageRespond(w, err.Error(), http.StatusOK)
				return
			}
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(items, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPost:
		var item entities.InventoryItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.InventoryService.CreateInventoryItem(item)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrInventoryItemAlreadyExists:
				statusCode = http.StatusConflict
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}
		jsonMessageRespond(w, fmt.Sprintf("Successfully created Inventory item with id %s", item.IngredientID), http.StatusCreated)
		return
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /inventory/{id}
func HandleInventoryItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		item, err := serviceinstance.InventoryService.GetInventoryItem(id)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrInventoryItemDoesntExist:
				statusCode = http.StatusNotFound
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}

		jsonPayload, err := json.MarshalIndent(item, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPut:
		var item entities.InventoryItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.InventoryService.UpdateInventoryItem(id, item)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrInventoryItemDoesntExist:
				statusCode = http.StatusNotFound
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}
		jsonMessageRespond(w, "Inventory Item successfully updated", http.StatusOK)
		return
	case http.MethodDelete:
		err := serviceinstance.InventoryService.DeleteInventoryItem(id)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrInventoryItemDoesntExist:
				statusCode = http.StatusNotFound
			}
			jsonErrorRespond(w, err, statusCode)
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

// Route: /inventory/getLeftOvers?sortBy={value}&page={page}&pageSize={pageSize}
func HandleInventoryLeftovers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sortBy := r.URL.Query().Get("sortBy")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil && pageStr != "" {
		jsonErrorRespond(w, ErrNonIntegerPage, http.StatusBadRequest)
		return
	}

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil && pageSizeStr != "" {
		jsonErrorRespond(w, ErrNonIntegerPageSize, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.InventoryService.GetLeftovers(sortBy, page, pageSize)
		if err != nil {
			statusCode := http.StatusInternalServerError
			switch {
			case errors.Is(err, serviceinstance.ErrInvalidSortValue),
				errors.Is(err, serviceinstance.ErrNegativePage):
				statusCode = http.StatusBadRequest
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}

		jsonPayload, err := json.MarshalIndent(items, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	default:
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
