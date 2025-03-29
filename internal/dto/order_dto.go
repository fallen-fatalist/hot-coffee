package dto

type OrderReport struct {
	ID           int64   `json:"order_id"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status,omitempty"`
	Total        float64 `json:"total,omitempty"`
	Reason       string  `json:"reason,omitempty"`
}
