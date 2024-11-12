package serviceinstance

import (
	"hot-coffee/internal/repositories/jsonrepository"
	"log/slog"
)

// Services instances
var (
	InventoryService *inventoryService
	MenuService      *menuService
	OrderService     *orderService
)

// Initialize services
func Init() {
	InventoryService = NewInventoryService(jsonrepository.NewInventoryRepository())
	MenuService = NewMenuService(jsonrepository.NewMenuRepository())
	// OrderService = NewOrderService()
	slog.Info("Services initialized")
}
