package entities

// TODO: Convert all IDs into int64
type MenuItem struct {
	ID          string               `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}

type MenuItemIngredient struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

type MenuItemSales struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	SalesCount  int    `json:"total_sales_count"`
}

type MenuItemSalesByCount []MenuItemSales

func (m MenuItemSalesByCount) Len() int {
	return len(m)
}

func (m MenuItemSalesByCount) Less(i, j int) bool {
	return m[i].SalesCount > m[j].SalesCount
}

func (m MenuItemSalesByCount) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type MenuReport struct {
	ID          string  `json:"product_id,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price,omitempty"`
	Relevance   float64 `json:"relevance"`
}
