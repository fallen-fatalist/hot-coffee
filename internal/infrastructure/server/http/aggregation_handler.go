package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"hot-coffee/internal/service/serviceinstance"
)

// Errors
var (
	ErrNonIntegerYear = fmt.Errorf("year must be an integer")
)

// Route: GET /reports/total-sales
func HandleTotalSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	sales, err := serviceinstance.OrderService.GetTotalSales()
	if err != nil {
		jsonErrorRespond(w, err, http.StatusBadRequest)
		return
	}

	jsonPayload, err := json.MarshalIndent(sales, "", "   ")
	if err != nil {
		jsonErrorRespond(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(jsonPayload)
}

// Route: GET /reports/popular-items
func HandlePopularItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	popularItems, err := serviceinstance.OrderService.GetPopularMenuItems()
	if err != nil {
		jsonErrorRespond(w, err, http.StatusInternalServerError)
		return
	}

	jsonPayload, err := json.MarshalIndent(popularItems, "", "   ")
	if err != nil {
		jsonErrorRespond(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(jsonPayload)
}

// Route: GET /reports/orderedItemsByPeriod
func HandleOrderedItemsByPeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	period := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	year, err := strconv.Atoi(yearStr)
	if err != nil && yearStr != "" {
		jsonErrorRespond(w, ErrNonIntegerYear, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.OrderService.GetOrderedItemsByPeriod(period, month, year)
		if err != nil {
			statusCode := http.StatusInternalServerError
			switch {
			case errors.Is(err, serviceinstance.ErrPeriodDayInvalid),
				errors.Is(err, serviceinstance.ErrPeriodTypeInvalid),
				errors.Is(err, serviceinstance.ErrPeriodMonthInvalid),
				errors.Is(err, serviceinstance.ErrParameterInvalid):
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

// Route: GET /orders/numberOfOrderedItems?startDate={startDate}&endDate={endDate}
func HandleNumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	switch r.Method {
	case http.MethodGet:
		items, err := serviceinstance.OrderService.GetOrderedMenuItemsCountByPeriod(startDate, endDate)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, serviceinstance.ErrEndDateEarlierThanStartDate) || errors.Is(err, serviceinstance.ErrInvalidDate) {
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

// Route: GET /reports/search?q=chocolate cake&filter=menu,orders&minPrice=10
func HandleFullTextSearchReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryString := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	switch r.Method {
	case http.MethodGet:
		orders, err := serviceinstance.AggregationService.FullTextSearchReport(queryString, filter, minPrice, maxPrice)
		if err != nil {
			statusCode := http.StatusInternalServerError
			switch {
			case errors.Is(err, serviceinstance.ErrNonNumericMaxPrice),
				errors.Is(err, serviceinstance.ErrNonNumericMinPrice),
				errors.Is(err, serviceinstance.ErrInvalidArgInFilter),
				errors.Is(err, serviceinstance.ErrMinMoreThanMaxPrice),
				errors.Is(err, serviceinstance.ErrNegativeMaxPrice),
				errors.Is(err, serviceinstance.ErrNegativeMinPrice),
				errors.Is(err, serviceinstance.ErrTooMuchArgs),
				errors.Is(err, serviceinstance.ErrMissingQueryString):
				statusCode = http.StatusBadRequest
			}
			jsonErrorRespond(w, err, statusCode)
			return
		}

		jsonPayload, err := json.MarshalIndent(orders, "", "   ")
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
