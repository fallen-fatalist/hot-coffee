package entities

type Order struct {
	ID           string      `json:"order_id,omitempty"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID         string `json:"product_id"`
	Quantity          int    `json:"quantity"`
	CustomizationInfo string `json:"customization_info"`
}

type TotalSales struct {
	Total float64 `json:"total_sales"`
}

type OrderedItemsCountByPeriod struct {
	Period            string         `json:"period"`
	Month             string         `json:"month,omitempty"`
	Year              int            `json:"year,omitempty"`
	OrderedItemsCount map[string]int `json:"orderedItems"`
}
