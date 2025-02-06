package vo

// view object

type InventoryUpdate struct {
	IngredientID string  `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity_used"`
	Remaining    float64 `json:"remaining"`
}
