package controllers

import (
	"encoding/json"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/services/serviceinstance"
	"hot-coffee/internal/utils"
	"net/http"
)

// Route: /menu
func HandleMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.MenuService.GetMenuItems()
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
		var item entities.MenuItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			utils.JSONErrorRespond(w, err)
			return
		}
		err = serviceinstance.MenuService.CreateMenuItem(item)
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

// Route: /menu/<id>
func HandleMenuItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
