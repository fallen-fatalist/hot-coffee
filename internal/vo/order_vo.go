package vo

import "hot-coffee/internal/dto"

// view object

type Response struct {
	Message string `json:"message"`
}

type BatchResponse struct {
	OrderReports []dto.OrderReport `json:"processed_orders"`
	Summary      Summary           `json:"summary"`
}

type Summary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted,omitempty"`
	Rejected         int               `json:"rejected,omitempty"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}
