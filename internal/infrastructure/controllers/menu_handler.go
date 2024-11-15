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

// Route: /menu
func HandleMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.MenuService.GetMenuItems()
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}

		jsonPayload, err := json.MarshalIndent(items, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPost:
		var item entities.MenuItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.MenuService.CreateMenuItem(item)
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

// Route: /menu/<id>
func HandleMenuItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	id := r.PathValue("id")
	// ID validation
	switch r.Method {
	case http.MethodGet:
		item, err := serviceinstance.MenuService.GetMenuItem(id)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
		}

		jsonPayload, err := json.MarshalIndent(item, "", "   ")
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		w.Write(jsonPayload)
		return
	case http.MethodPut:
		var item entities.MenuItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.MenuService.UpdateMenuItem(id, item)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
			return
		}
		return
	case http.MethodDelete:
		err := serviceinstance.MenuService.DeleteMenuItem(id)
		if err != nil {
			utils.JSONErrorRespond(w, err, http.StatusBadRequest)
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
