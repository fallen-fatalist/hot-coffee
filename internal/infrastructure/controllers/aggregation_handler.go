package controllers

import "net/http"

// Route: /reports/total-sales
func HandleTotalSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

// Route: /reports/popular-items
func HandlePopularItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}
