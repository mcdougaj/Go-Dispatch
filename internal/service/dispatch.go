package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mcdougaj/Go-Dispatch/internal/api/googlemaps"
	"github.com/mcdougaj/Go-Dispatch/internal/api/motive"
	"github.com/mcdougaj/Go-Dispatch/internal/database"
	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// DispatchService handles dispatch operations
type DispatchService struct {
	googleMaps   *googlemaps.Client
	motive       *motive.Client
	driverRepo   *database.DriverRepository
	vehicleRepo  *database.VehicleRepository
	routeRepo    *database.RouteRepository
}

// NewDispatchService creates a new dispatch service
func NewDispatchService(
	googleMaps *googlemaps.Client,
	motive *motive.Client,
	db *database.DB,
) *DispatchService {
	return &DispatchService{
		googleMaps:  googleMaps,
		motive:      motive,
		driverRepo:  database.NewDriverRepository(db),
		vehicleRepo: database.NewVehicleRepository(db),
		routeRepo:   database.NewRouteRepository(db),
	}
}

// CreateRoute creates a new dispatch route (saved to local database)
func (s *DispatchService) CreateRoute(ctx context.Context, req models.DispatchRequest) (*models.Route, error) {
	// Validate driver and vehicle exist in our database
	if _, err := s.driverRepo.GetByID(ctx, req.DriverID); err != nil {
		return nil, fmt.Errorf("driver not found: %s", req.DriverID)
	}
	if _, err := s.vehicleRepo.GetByID(ctx, req.VehicleID); err != nil {
		return nil, fmt.Errorf("vehicle not found: %s", req.VehicleID)
	}

	// Generate IDs for stops
	for i := range req.Stops {
		if req.Stops[i].ID == "" {
			req.Stops[i].ID = uuid.New().String()
		}
		if req.Stops[i].Status == "" {
			req.Stops[i].Status = "pending"
		}
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

	// Calculate route metrics using Google Maps
	if err := s.calculateRouteMetrics(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to calculate route metrics: %w", err)
	}

	// Save to database
	if err := s.routeRepo.Create(ctx, route); err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	return route, nil
}

// OptimizeRoute optimizes the order of stops in a route using Google Maps
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

// GetRoute retrieves a route by ID from database
func (s *DispatchService) GetRoute(ctx context.Context, routeID string) (*models.Route, error) {
	return s.routeRepo.GetByID(ctx, routeID)
}

// ListRoutes retrieves all routes from database
func (s *DispatchService) ListRoutes(ctx context.Context) ([]*models.Route, error) {
	return s.routeRepo.List(ctx)
}

// UpdateRouteStatus updates the status of a route in database
func (s *DispatchService) UpdateRouteStatus(ctx context.Context, routeID, status string) error {
	// Get existing route to check current state
	route, err := s.routeRepo.GetByID(ctx, routeID)
	if err != nil {
		return err
	}

	var startTime, endTime sql.NullTime
	if route.StartTime != nil {
		startTime = sql.NullTime{Time: *route.StartTime, Valid: true}
	}
	if route.EndTime != nil {
		endTime = sql.NullTime{Time: *route.EndTime, Valid: true}
	}

	if status == "in_progress" && !startTime.Valid {
		startTime = sql.NullTime{Time: time.Now(), Valid: true}
	} else if status == "completed" && !endTime.Valid {
		endTime = sql.NullTime{Time: time.Now(), Valid: true}
	}

	return s.routeRepo.UpdateStatus(ctx, routeID, status, &startTime, &endTime)
}

// SyncVehiclesFromMotive syncs vehicle data from Motive (READ-ONLY from Motive, saves to local DB)
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

		// Save to local database (upsert by motive_id)
		if err := s.vehicleRepo.CreateFromMotive(ctx, vehicle, mv.ID); err != nil {
			return fmt.Errorf("failed to save vehicle from Motive: %w", err)
		}
	}

	return nil
}

// SyncDriversFromMotive syncs driver data from Motive (READ-ONLY from Motive, saves to local DB)
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

		// Save to local database (upsert by motive_id)
		if err := s.driverRepo.CreateFromMotive(ctx, driver, md.ID); err != nil {
			return fmt.Errorf("failed to save driver from Motive: %w", err)
		}
	}

	return nil
}

// GetVehicleLocation gets the current location of a vehicle from Motive and updates local DB
func (s *DispatchService) GetVehicleLocation(ctx context.Context, vehicleID string) (*models.Location, error) {
	vehicle, err := s.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}

	if vehicle.MotiveID == "" {
		return nil, fmt.Errorf("vehicle not linked to Motive")
	}

	// Get live location from Motive
	location, err := s.motive.GetVehicleLocation(ctx, vehicle.MotiveID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle location from Motive: %w", err)
	}

	// Update location in local database
	vehicle.Location = location
	vehicle.UpdatedAt = time.Now()
	if err := s.vehicleRepo.Update(ctx, vehicle); err != nil {
		return nil, fmt.Errorf("failed to update vehicle location: %w", err)
	}

	return location, nil
}

// calculateRouteMetrics calculates distance and duration for a route using Google Maps
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

// CreateDriver creates a new driver in local database
func (s *DispatchService) CreateDriver(ctx context.Context, driver *models.Driver) error {
	driver.ID = uuid.New().String()
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()
	if driver.Status == "" {
		driver.Status = "available"
	}
	return s.driverRepo.Create(ctx, driver)
}

// CreateVehicle creates a new vehicle in local database
func (s *DispatchService) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	vehicle.ID = uuid.New().String()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()
	if vehicle.Status == "" {
		vehicle.Status = "active"
	}
	return s.vehicleRepo.Create(ctx, vehicle)
}

// ListDrivers retrieves all drivers from database
func (s *DispatchService) ListDrivers(ctx context.Context) ([]*models.Driver, error) {
	return s.driverRepo.List(ctx)
}

// ListVehicles retrieves all vehicles from database
func (s *DispatchService) ListVehicles(ctx context.Context) ([]*models.Vehicle, error) {
	return s.vehicleRepo.List(ctx)
}
