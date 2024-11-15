package controllers

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/services/serviceinstance"
	"hot-coffee/internal/utils"
)

// Route: /reports/total-sales
func HandleTotalSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sales, err := serviceinstance.OrderService.GetTotalSales()
	if err != nil {
		utils.JSONErrorRespond(w, err, http.StatusBadRequest)
		return
	}

	jsonPayload, err := json.MarshalIndent(sales, "", "   ")
	if err != nil {
		utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(jsonPayload)
	return
}

// Route: /reports/popular-items
func HandlePopularItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	popularItems, err := serviceinstance.OrderService.GetPopularMenuItems()
	if err != nil {
		utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
		return
	}

	jsonPayload, err := json.MarshalIndent(popularItems, "", "   ")
	if err != nil {
		utils.JSONErrorRespond(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(jsonPayload)
	return
}
