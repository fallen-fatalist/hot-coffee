package entities

// TODO: Convert all IDs into int64
type Order struct {
	ID           string      `json:"order_id,omitempty"`
	CustomerName string      `json:"customer_name"`
	CustomerID   int64       `json:"customer_id,omitempty"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status,omitempty"`
	CreatedAt    string      `json:"created_at,omitempty"`
}

type OrderItem struct {
	ProductID         int    `json:"product_id"`
	Quantity          int    `json:"quantity"`
	CustomizationInfo string `json:"customization_info,omitempty"`
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

type OrderedMenuItemsCount struct {
	Espresso            int `json:"espresso,omitempty"`
	Latte               int `json:"latte,omitempty"`
	Cappuccino          int `json:"cappuccino,omitempty"`
	Americano           int `json:"americano,omitempty"`
	FlatWhite           int `json:"flat_white,omitempty"`
	Mocha               int `json:"mocha,omitempty"`
	Croissant           int `json:"croissant,omitempty"`
	Muffin              int `json:"muffin,omitempty"`
	BlueberryMuffin     int `json:"blueberry_muffin,omitempty"`
	ChocolateChipCookie int `json:"chocolate_chip_cookie,omitempty"`
	Bagel               int `json:"bagel,omitempty"`
	Cheesecake          int `json:"cheesecake,omitempty"`
	Tiramisu            int `json:"tiramisu,omitempty"`
	ChocolateCake       int `json:"chocolate_cake,omitempty"`
	VanillaCupcake      int `json:"vanilla_cupcake,omitempty"`
}

type OrderReport struct {
	ID           string   `json:"order_id,omitempty"`
	CustomerName string   `json:"customer_name"`
	Items        []string `json:"items"`
	Total        float64  `json:"total,omitempty"`
	Relevance    float64  `json:"relevance"`
}
