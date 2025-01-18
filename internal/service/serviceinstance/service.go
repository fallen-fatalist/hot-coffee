package serviceinstance

import (
	"hot-coffee/internal/infrastructure/storage/jsonrepository"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"log/slog"
)

// Services instances
var (
	InventoryService service.InventoryService
	MenuService      service.MenuService
	OrderService     service.OrderService
)

func NewService(repositories *repository.Repository) (*service.Service, error) {
	inventoryService, err := NewInventoryService(repositories.Inventory)
	if err != nil {
		return nil, err
	}
	return &service.Service{
		InventoryService: inventoryService,
		MenuService:      NewMenuService(repositories.Menu),
		OrderService:     NewOrderService(repositories.Order),
	}, nil
}

// Initialize services
func Init() {
	var err error
	serviceInstance, err := NewService(jsonrepository.NewRepository())
	utils.FatalError("Error while initializing inventory service", err)
	InventoryService = serviceInstance.InventoryService
	MenuService = serviceInstance.MenuService
	OrderService = serviceInstance.OrderService
	slog.Info("Services initialized")
}
