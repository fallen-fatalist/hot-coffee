package vo

import "hot-coffee/internal/dto"

// view object

type BatchResponse struct {
	ProcessedOrders []dto.ProcessedOrder `json:"processed_orders"`
	Summary         Summary              `json:"summary"`
}

type Summary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted,omitempty"`
	Rejected         int               `json:"rejected,omitempty"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}
