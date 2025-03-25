package entities

type FullReport struct {
	Menus        []MenuReport  `json:"menu_items,omitempty"`
	Orders       []OrderReport `json:"orders,omitempty"`
	TotalMatches int           `json:"total_matches"`
}
