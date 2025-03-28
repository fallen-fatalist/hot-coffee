package serviceinstance

import (
	"hot-coffee/internal/infrastructure/storage/postgres"
	"hot-coffee/internal/repository"
	"hot-coffee/internal/service"
	"hot-coffee/internal/utils"
	"log/slog"
)

// Services instances
var (
	InventoryService   service.InventoryService
	MenuService        service.MenuService
	OrderService       service.OrderService
	AggregationService service.AggregationService // New aggregation service
)

func NewService(repositories *repository.Repository) (*service.Service, error) {

	return &service.Service{
		InventoryService:   NewInventoryService(repositories.Inventory),
		MenuService:        NewMenuService(repositories.Menu),
		OrderService:       NewOrderService(repositories.Order),
		AggregationService: NewAggregationService(repositories.Menu, repositories.Order), // New aggregation service
	}, nil
}

// Initialize services
func Init() {
	var err error
	serviceInstance, err := NewService(postgres.NewRepository())
	utils.FatalError("Error while initializing inventory service", err)
	InventoryService = serviceInstance.InventoryService
	MenuService = serviceInstance.MenuService
	OrderService = serviceInstance.OrderService
	AggregationService = serviceInstance.AggregationService // New aggregation service
	slog.Info("Services initialized")
}
