package controller

import (
	"net/http"
)

// Route: /orders
func HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//service.OrderService.GetOrder(w, r)
	case http.MethodPost:
		//service.OrderService.PostOrderItem(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orders/<id>
func HandleOrder(w http.ResponseWriter, r *http.Request) {
	// ID validation
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /orderds/<id>/close
func HandleOrderClose(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//service.OrderService.CloseOrder(w, r)
	} else {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
