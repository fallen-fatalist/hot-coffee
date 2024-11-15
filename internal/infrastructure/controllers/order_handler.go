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

// Route: /orders
func HandleOrders(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		orders, err := serviceinstance.OrderService.GetOrders()
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(orders, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPost:
		var order entities.Order
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.OrderService.CreateOrder(order)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orders/<id>
func HandleOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		order, err := serviceinstance.OrderService.GetOrder(id)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
		}
		jsonPayload, err := json.MarshalIndent(order, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPut:
		var order entities.Order
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&order)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.OrderService.UpdateOrder(id, order)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		return
	case http.MethodDelete:
		err := serviceinstance.OrderService.DeleteOrder(id)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orderds/<id>/close
func HandleOrderClose(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	id := r.PathValue("id")
	if r.Method == http.MethodPost {
		serviceinstance.OrderService.CloseOrder(id)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func HandleOpenOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		openOrders, err := serviceinstance.OrderService.GetOpenOrders()
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(openOrders, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	} else {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
