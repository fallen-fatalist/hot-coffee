package entities

type InventoryItem struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type PaginatedInventoryItems struct {
	CurrentPage int                 `json:"currentPage"`
	HasNextPage bool                `json:"hasNextPage"`
	PageSize    int                 `json:"pageSize"`
	TotalPages  int                 `json:"totalPages"`
	Items       []PageInventoryItem `json:"data"`
}

type PageInventoryItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}
