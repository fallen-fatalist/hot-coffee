package serviceinstance

import (
	"log/slog"

	"hot-coffee/internal/repositories/jsonrepository"
	"hot-coffee/internal/services/service"
)

// Services instances
var (
	InventoryService service.InventoryService
	MenuService      service.MenuService
	OrderService     service.OrderService
)

// Initialize services
func Init() {
	InventoryService = NewInventoryService(jsonrepository.NewInventoryRepository())
	MenuService = NewMenuService(jsonrepository.NewMenuRepository())
	// OrderService = NewOrderService()
	slog.Info("Services initialized")
}
