package serviceinstance

import (
	"errors"
	"hot-coffee/internal/core/entities"
	"hot-coffee/internal/repository"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type aggService struct {
	menuRepository  repository.MenuRepository
	orderRepository repository.OrderRepository
}

// Errors
var (
	ErrInvalidArgInFilter  = errors.New("invalid arguments in filter field")
	ErrTooMuchArgs         = errors.New("too much arguments in filter field")
	ErrNonNumericMinPrice  = errors.New("non-numeric minPrice is provided")
	ErrNegativeMinPrice    = errors.New("negative minPrice is provided")
	ErrNonNumericMaxPrice  = errors.New("non-numeric maxPrice is provided")
	ErrNegativeMaxPrice    = errors.New("negative maxPrice is provided")
	ErrMinMoreThanMaxPrice = errors.New("minPrice cannot be more than or equal to maxPrice")
)

// Valid filter options
var validFilters = map[string]struct{}{
	"menu":   {},
	"orders": {},
	"all":    {},
}

func NewAggregationService(menuRepo repository.MenuRepository, orderRepo repository.OrderRepository) *aggService {
	if menuRepo == nil || orderRepo == nil {
		slog.Error("Error while creating Aggregation service: Nil pointer repository provided")
		os.Exit(1)
	}

	return &aggService{
		menuRepository:  menuRepo,
		orderRepository: orderRepo,
	}
}

func (s *aggService) FullTextSearchReport(q, filter, minPriceStr, maxPriceStr string) (entities.FullReport, error) {
	var result entities.FullReport

	// Convert price strings to integers
	minPrice, err := parsePrice(minPriceStr, ErrNonNumericMinPrice, ErrNegativeMinPrice)
	if err != nil {
		return result, err
	}

	maxPrice, err := parsePrice(maxPriceStr, ErrNonNumericMaxPrice, ErrNegativeMaxPrice)
	if err != nil {
		return result, err
	}

	// Ensure minPrice is not greater than maxPrice (unless maxPrice is 0, meaning no max limit)
	if minPrice > 0 && maxPrice > 0 && minPrice >= maxPrice {
		return result, ErrMinMoreThanMaxPrice
	}

	if filter == "" {
		return s.fetchBothReports(q, minPrice, maxPrice)
	}

	// Parse and validate filter options
	filterData := strings.Split(filter, ",")
	options := make(map[string]bool)

	for _, data := range filterData {
		lowData := strings.ToLower(strings.TrimSpace(data))
		if _, valid := validFilters[lowData]; !valid {
			return result, ErrInvalidArgInFilter
		}
		options[lowData] = true
	}

	// Determine which reports to fetch
	isMenu := options["menu"]
	isOrders := options["orders"]
	isAll := options["all"]

	if isAll && (isMenu || isOrders) {
		return result, ErrTooMuchArgs
	}

	if isAll || (isMenu && isOrders) {
		return s.fetchBothReports(q, minPrice, maxPrice)
	} else if isMenu {
		return s.fetchMenuReport(q, minPrice, maxPrice)
	} else if isOrders {
		return s.fetchOrdersReport(q, minPrice, maxPrice)
	}

	return result, nil
}

// Fetch both menu and order reports
func (s *aggService) fetchBothReports(q string, minPrice, maxPrice int) (entities.FullReport, error) {
	menus, err := s.menuRepository.GetMenusFullTextSearchReport(q, minPrice, maxPrice)
	if err != nil {
		return entities.FullReport{}, err
	}

	orders, err := s.orderRepository.GetOrdersFullTextSearchReport(q, minPrice, maxPrice)
	if err != nil {
		return entities.FullReport{}, err
	}

	return entities.FullReport{
		Menus:        menus,
		Orders:       orders,
		TotalMatches: len(menus) + len(orders),
	}, nil
}

// Fetch only menu report
func (s *aggService) fetchMenuReport(q string, minPrice, maxPrice int) (entities.FullReport, error) {
	menus, err := s.menuRepository.GetMenusFullTextSearchReport(q, minPrice, maxPrice)
	if err != nil {
		return entities.FullReport{}, err
	}

	return entities.FullReport{
		Menus:        menus,
		TotalMatches: len(menus),
	}, nil
}

// Fetch only order report
func (s *aggService) fetchOrdersReport(q string, minPrice, maxPrice int) (entities.FullReport, error) {
	orders, err := s.orderRepository.GetOrdersFullTextSearchReport(q, minPrice, maxPrice)
	if err != nil {
		return entities.FullReport{}, err
	}

	return entities.FullReport{
		Orders:       orders,
		TotalMatches: len(orders),
	}, nil
}

// Helper function to parse price values
func parsePrice(priceStr string, errNonNumeric, errNegative error) (int, error) {
	if priceStr == "" {
		return 0, nil
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0, errNonNumeric
	}

	if price < 0 {
		return 0, errNegative
	}

	return price, nil
}
