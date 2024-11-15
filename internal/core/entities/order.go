package entities

type Order struct {
	ID           string      `json:"order_id,omitempty"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type TotalSales struct {
	Total float64 `json:"total_sales"`
}
