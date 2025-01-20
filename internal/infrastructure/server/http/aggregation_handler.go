package httpsever

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"hot-coffee/internal/service/serviceinstance"
	"hot-coffee/internal/utils"
)

// Errors
var (
	ErrNonIntegerYear = fmt.Errorf("year must be an integer")
)

// Route: /reports/total-sales
func HandleTotalSales(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
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
}

// Route: /reports/popular-items
func HandlePopularItems(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
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
}

// Route: /reports/orderedItemsByPeriod
func HandleOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info(fmt.Sprintf("%s request with URL: %s", r.Method, r.URL.String()))

	period := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil && yearStr != "" {
		utils.JSONErrorRespond(w, ErrNonIntegerYear, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.OrderService.GetOrderedItemsByPeriod(period, month, year)
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
	default:
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
