package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"hot-coffee/internal/service/serviceinstance"
)

// Route: /orders
func HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		orders, err := serviceinstance.OrderService.GetOrders()
		if err != nil {
			if errors.Is(err, serviceinstance.ErrNoOrders) {
				jsonMessageRespond(w, "No orders", http.StatusOK)
				return
			}
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(orders, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var order entities.Order
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		orderID, err := serviceinstance.OrderService.CreateOrder(order)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusBadRequest)
			return
		}

		jsonMessageRespond(w, fmt.Sprintf("Successfully created Order with ID: %d", orderID), http.StatusCreated)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orders/<id>

// MUST DO: Handle no rows in result set, not found IDs for orders, menu and inventory items
func HandleOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		order, err := serviceinstance.OrderService.GetOrder(id)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		jsonPayload, err := json.MarshalIndent(order, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPut:
		var order entities.Order
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.OrderService.UpdateOrder(id, order)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		return
	case http.MethodDelete:
		err := serviceinstance.OrderService.DeleteOrder(id)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orderds/<id>/close
func HandleOrderClose(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	switch r.Method {
	case http.MethodPost:
		err := serviceinstance.OrderService.CloseOrder(id)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func HandleOrderInProgress(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	switch r.Method {
	case http.MethodPost:
		err := serviceinstance.OrderService.SetInProgress(id)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "POST, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func HandleOpenOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		openOrders, err := serviceinstance.OrderService.GetOpenOrders()
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(openOrders, "", "   ")
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "GET, OPTIONS")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type batchRequest struct {
	Orders []entities.Order `json:"orders"`
}

// POST /orders/batch-process
func HandleBatchOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		// Request Body  Parsing \\
		var req batchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonErrorRespond(w, errors.ErrIncorrectRequest, http.StatusBadRequest)
			return
		}

		// Service Call \\
		response, err := serviceinstance.OrderService.CreateOrders(req.Orders)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		json, err := json.Marshal(response)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		// Response \\
		w.WriteHeader(http.StatusCreated)
		w.Write(json)
		return
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}
