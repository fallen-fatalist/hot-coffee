package httpserver

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/core/errors"
	"hot-coffee/internal/infrastructure/storage/jsonrepository"
	"hot-coffee/internal/service/serviceinstance"
)

// Route: /menu
func HandleMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.MenuService.GetMenuItems()
		if err != nil {
			if errors.Is(err, serviceinstance.ErrNoMenuItems) {
				jsonMessageRespond(w, "no menu items", http.StatusOK)
			}
			statusCode := http.StatusBadRequest
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
	case http.MethodPost:
		var item entities.MenuItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.MenuService.CreateMenuItem(item)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrMenuItemAlreadyExists:
				statusCode = http.StatusConflict
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}
		jsonMessageRespond(w, "Menu Item successfully created", http.StatusCreated)
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

	id := r.PathValue("id")
	switch r.Method {
	case http.MethodGet:
		item, err := serviceinstance.MenuService.GetMenuItem(id)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrMenuItemNotExists:
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
		var item entities.MenuItem
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&item)
		if err != nil {
			jsonErrorRespond(w, err, http.StatusInternalServerError)
			return
		}
		err = serviceinstance.MenuService.UpdateMenuItem(id, item)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case serviceinstance.ErrMenuItemNotExists:
				statusCode = http.StatusNotFound
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}
		return
	case http.MethodDelete:
		err := serviceinstance.MenuService.DeleteMenuItem(id)
		if err != nil {
			statusCode := http.StatusBadRequest
			switch err {
			case jsonrepository.ErrMenuItemDoesntExist:
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
