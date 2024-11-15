package serviceinstance

import (
	"log/slog"

	"hot-coffee/internal/repositories/jsonrepository"
	"hot-coffee/internal/services/service"
	"hot-coffee/internal/utils"
)

// Services instances
var (
	InventoryService service.InventoryService
	MenuService      service.MenuService
	OrderService     service.OrderService
)

// Initialize services
func Init() {
	var err error
	InventoryService, err = NewInventoryService(jsonrepository.NewInventoryRepository())
	utils.FatalError("Error while initializing inventory service", err)
	MenuService = NewMenuService(jsonrepository.NewMenuRepository())
	OrderService = NewOrderService(jsonrepository.NewOrderRepository())
	slog.Info("Services initialized")
}
