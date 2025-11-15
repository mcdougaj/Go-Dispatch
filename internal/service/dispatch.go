package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mcdougaj/Go-Dispatch/internal/api/googlemaps"
	"github.com/mcdougaj/Go-Dispatch/internal/api/motive"
	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// DispatchService handles dispatch operations
type DispatchService struct {
	googleMaps *googlemaps.Client
	motive     *motive.Client
	routes     map[string]*models.Route // In-memory storage (would use DB in production)
	drivers    map[string]*models.Driver
	vehicles   map[string]*models.Vehicle
}

// NewDispatchService creates a new dispatch service
func NewDispatchService(googleMaps *googlemaps.Client, motive *motive.Client) *DispatchService {
	return &DispatchService{
		googleMaps: googleMaps,
		motive:     motive,
		routes:     make(map[string]*models.Route),
		drivers:    make(map[string]*models.Driver),
		vehicles:   make(map[string]*models.Vehicle),
	}
}

// CreateRoute creates a new dispatch route
func (s *DispatchService) CreateRoute(ctx context.Context, req models.DispatchRequest) (*models.Route, error) {
	// Validate driver and vehicle
	if _, exists := s.drivers[req.DriverID]; !exists {
		return nil, fmt.Errorf("driver not found: %s", req.DriverID)
	}
	if _, exists := s.vehicles[req.VehicleID]; !exists {
		return nil, fmt.Errorf("vehicle not found: %s", req.VehicleID)
	}

	// Create the route
	route := &models.Route{
		ID:        uuid.New().String(),
		DriverID:  req.DriverID,
		VehicleID: req.VehicleID,
		Stops:     req.Stops,
		Status:    "planned",
		StartTime: req.StartTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate route metrics
	if err := s.calculateRouteMetrics(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to calculate route metrics: %w", err)
	}

	// Store the route
	s.routes[route.ID] = route

	return route, nil
}

// OptimizeRoute optimizes the order of stops in a route
func (s *DispatchService) OptimizeRoute(ctx context.Context, req models.RouteOptimizationRequest) (*models.RouteOptimizationResponse, error) {
	if len(req.Destinations) == 0 {
		return nil, fmt.Errorf("no destinations provided")
	}

	// Use Google Maps to optimize waypoint order
	optimizedOrder, err := s.googleMaps.OptimizeWaypoints(ctx, req.Origin, req.Destinations)
	if err != nil {
		return nil, fmt.Errorf("failed to optimize waypoints: %w", err)
	}

	// Calculate total distance and duration
	var totalDistance float64
	var totalDuration int

	for i := 0; i < len(req.Destinations); i++ {
		var origin models.Location
		if i == 0 {
			origin = req.Origin
		} else {
			origin = req.Destinations[optimizedOrder[i-1]]
		}

		destination := req.Destinations[optimizedOrder[i]]
		routeInfo, err := s.googleMaps.CalculateRoute(ctx, origin, destination)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate route segment: %w", err)
		}

		totalDistance += float64(routeInfo.DistanceMeters) / 1609.34 // Convert meters to miles
		totalDuration += int(routeInfo.DurationSeconds / 60)         // Convert seconds to minutes
	}

	return &models.RouteOptimizationResponse{
		OptimizedOrder: optimizedOrder,
		TotalDistance:  totalDistance,
		TotalDuration:  totalDuration,
	}, nil
}

// GetRoute retrieves a route by ID
func (s *DispatchService) GetRoute(ctx context.Context, routeID string) (*models.Route, error) {
	route, exists := s.routes[routeID]
	if !exists {
		return nil, fmt.Errorf("route not found: %s", routeID)
	}
	return route, nil
}

// ListRoutes retrieves all routes
func (s *DispatchService) ListRoutes(ctx context.Context) ([]*models.Route, error) {
	routes := make([]*models.Route, 0, len(s.routes))
	for _, route := range s.routes {
		routes = append(routes, route)
	}
	return routes, nil
}

// UpdateRouteStatus updates the status of a route
func (s *DispatchService) UpdateRouteStatus(ctx context.Context, routeID, status string) error {
	route, exists := s.routes[routeID]
	if !exists {
		return fmt.Errorf("route not found: %s", routeID)
	}

	route.Status = status
	route.UpdatedAt = time.Now()

	if status == "in_progress" && route.StartTime == nil {
		now := time.Now()
		route.StartTime = &now
	} else if status == "completed" && route.EndTime == nil {
		now := time.Now()
		route.EndTime = &now
	}

	return nil
}

// SyncVehiclesFromMotive syncs vehicle data from Motive
func (s *DispatchService) SyncVehiclesFromMotive(ctx context.Context) error {
	motiveVehicles, err := s.motive.GetVehicles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get vehicles from Motive: %w", err)
	}

	for _, mv := range motiveVehicles {
		vehicle := &models.Vehicle{
			ID:           uuid.New().String(),
			VIN:          mv.VIN,
			LicensePlate: mv.LicensePlate,
			Make:         mv.Make,
			Model:        mv.Model,
			Year:         mv.Year,
			Status:       "active",
			MotiveID:     mv.ID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Get current location if available
		if mv.CurrentLocation != nil {
			vehicle.Location = &models.Location{
				Latitude:  mv.CurrentLocation.Latitude,
				Longitude: mv.CurrentLocation.Longitude,
			}
		}

		s.vehicles[vehicle.ID] = vehicle
	}

	return nil
}

// SyncDriversFromMotive syncs driver data from Motive
func (s *DispatchService) SyncDriversFromMotive(ctx context.Context) error {
	motiveDrivers, err := s.motive.GetDrivers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get drivers from Motive: %w", err)
	}

	for _, md := range motiveDrivers {
		driver := &models.Driver{
			ID:          uuid.New().String(),
			Name:        fmt.Sprintf("%s %s", md.FirstName, md.LastName),
			PhoneNumber: md.PhoneNumber,
			Email:       md.Email,
			Status:      "available",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		s.drivers[driver.ID] = driver
	}

	return nil
}

// GetVehicleLocation gets the current location of a vehicle from Motive
func (s *DispatchService) GetVehicleLocation(ctx context.Context, vehicleID string) (*models.Location, error) {
	vehicle, exists := s.vehicles[vehicleID]
	if !exists {
		return nil, fmt.Errorf("vehicle not found: %s", vehicleID)
	}

	if vehicle.MotiveID == "" {
		return nil, fmt.Errorf("vehicle not linked to Motive")
	}

	location, err := s.motive.GetVehicleLocation(ctx, vehicle.MotiveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle location: %w", err)
	}

	// Update cached location
	vehicle.Location = location
	vehicle.UpdatedAt = time.Now()

	return location, nil
}

// calculateRouteMetrics calculates distance and duration for a route
func (s *DispatchService) calculateRouteMetrics(ctx context.Context, route *models.Route) error {
	if len(route.Stops) < 2 {
		return nil
	}

	var totalDistance float64
	var totalDuration int

	for i := 0; i < len(route.Stops)-1; i++ {
		origin := route.Stops[i].Location
		destination := route.Stops[i+1].Location

		routeInfo, err := s.googleMaps.CalculateRoute(ctx, origin, destination)
		if err != nil {
			return fmt.Errorf("failed to calculate route segment %d: %w", i, err)
		}

		totalDistance += float64(routeInfo.DistanceMeters) / 1609.34 // Convert to miles
		totalDuration += int(routeInfo.DurationSeconds / 60)         // Convert to minutes
	}

	route.TotalDistance = totalDistance
	route.EstimatedDuration = totalDuration

	return nil
}

// CreateDriver creates a new driver
func (s *DispatchService) CreateDriver(ctx context.Context, driver *models.Driver) error {
	driver.ID = uuid.New().String()
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()
	s.drivers[driver.ID] = driver
	return nil
}

// CreateVehicle creates a new vehicle
func (s *DispatchService) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	vehicle.ID = uuid.New().String()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	s.vehicles[vehicle.ID] = vehicle
	return nil
}

// ListDrivers retrieves all drivers
func (s *DispatchService) ListDrivers(ctx context.Context) ([]*models.Driver, error) {
	drivers := make([]*models.Driver, 0, len(s.drivers))
	for _, driver := range s.drivers {
		drivers = append(drivers, driver)
	}
	return drivers, nil
}

// ListVehicles retrieves all vehicles
func (s *DispatchService) ListVehicles(ctx context.Context) ([]*models.Vehicle, error) {
	vehicles := make([]*models.Vehicle, 0, len(s.vehicles))
	for _, vehicle := range s.vehicles {
		vehicles = append(vehicles, vehicle)
	}
	return vehicles, nil
}
