package controller

import "net/http"

// Route: /menu
func HandleMenu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Route: /menu/<id>
func HandleMenuItem(w http.ResponseWriter, r *http.Request) {
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
