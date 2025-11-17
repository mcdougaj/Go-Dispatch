package services

import (
	"context"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
	"github.com/mcdougaj/Go-Dispatch/internal/repositories"
	"googlemaps.github.io/maps"
)

type RouteService struct {
	routeRepo   *repositories.RouteRepository
	mapsClient  *maps.Client
}

func NewRouteService(routeRepo *repositories.RouteRepository, mapsClient *maps.Client) *RouteService {
	return &RouteService{
		routeRepo:   routeRepo,
		mapsClient:  mapsClient,
	}
}

func (s *RouteService) CreateRoute(req *models.CreateRouteRequest) (*models.Route, error) {
	// Create the route
	route, err := s.routeRepo.Create(req)
	if err != nil {
		return nil, err
	}

	// Optimize the route if we have stops
	if len(req.Stops) > 0 {
		if err := s.OptimizeRoute(route.ID); err != nil {
			// Log error but don't fail the creation
			fmt.Printf("Warning: failed to optimize route: %v\n", err)
		}
	}

	return route, nil
}

func (s *RouteService) GetRoute(id int) (*models.Route, error) {
	return s.routeRepo.GetByID(id)
}

func (s *RouteService) GetAllRoutes() ([]models.Route, error) {
	return s.routeRepo.GetAll()
}

func (s *RouteService) UpdateRouteStatus(id int, status string) error {
	return s.routeRepo.UpdateStatus(id, status)
}

func (s *RouteService) OptimizeRoute(routeID int) error {
	route, err := s.routeRepo.GetByID(routeID)
	if err != nil {
		return err
	}

	if len(route.Stops) < 2 {
		return nil // Nothing to optimize
	}

	// Build waypoints
	var waypoints []string
	for _, stop := range route.Stops {
		waypoints = append(waypoints, fmt.Sprintf("%f,%f", stop.Latitude, stop.Longitude))
	}

	// Get directions from Google Maps
	req := &maps.DirectionsRequest{
		Origin:      waypoints[0],
		Destination: waypoints[len(waypoints)-1],
		Mode:        maps.TravelModeDriving,
		Optimize:    true,
	}

	if len(waypoints) > 2 {
		req.Waypoints = waypoints[1 : len(waypoints)-1]
	}

	directions, _, err := s.mapsClient.Directions(context.Background(), req)
	if err != nil {
		return fmt.Errorf("failed to get directions: %w", err)
	}

	if len(directions) == 0 {
		return fmt.Errorf("no directions found")
	}

	// Calculate total distance and duration
	var totalDistance float64
	var totalDuration int

	for _, leg := range directions[0].Legs {
		totalDistance += float64(leg.Distance.Meters) / 1000.0 // Convert to kilometers
		totalDuration += int(leg.Duration.Minutes())
	}

	// Update route metrics
	return s.routeRepo.UpdateRouteMetrics(routeID, totalDistance, totalDuration)
}
