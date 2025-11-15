package googlemaps

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// Client wraps the Google Maps API client
type Client struct {
	client *maps.Client
}

// NewClient creates a new Google Maps API client
func NewClient(apiKey string) (*Client, error) {
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %w", err)
	}

	return &Client{client: c}, nil
}

// Geocode converts an address to coordinates
func (c *Client) Geocode(ctx context.Context, address string) (*models.Location, error) {
	req := &maps.GeocodingRequest{
		Address: address,
	}

	results, err := c.client.Geocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	location := &models.Location{
		Latitude:  results[0].Geometry.Location.Lat,
		Longitude: results[0].Geometry.Location.Lng,
		Address:   results[0].FormattedAddress,
	}

	return location, nil
}

// ReverseGeocode converts coordinates to an address
func (c *Client) ReverseGeocode(ctx context.Context, lat, lng float64) (string, error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
	}

	results, err := c.client.Geocode(ctx, req)
	if err != nil {
		return "", fmt.Errorf("reverse geocoding failed: %w", err)
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no results found for coordinates: %f, %f", lat, lng)
	}

	return results[0].FormattedAddress, nil
}

// CalculateRoute calculates a route between origin and destination
func (c *Client) CalculateRoute(ctx context.Context, origin, destination models.Location) (*RouteInfo, error) {
	req := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("%f,%f", origin.Latitude, origin.Longitude),
		Destination: fmt.Sprintf("%f,%f", destination.Latitude, destination.Longitude),
	}

	routes, _, err := c.client.Directions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("directions request failed: %w", err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	route := routes[0]
	var totalDistance int64
	var totalDuration int64

	for _, leg := range route.Legs {
		totalDistance += leg.Distance.Meters
		totalDuration += int64(leg.Duration.Seconds())
	}

	return &RouteInfo{
		DistanceMeters: totalDistance,
		DurationSeconds: totalDuration,
		Polyline:       route.OverviewPolyline.Points,
	}, nil
}

// OptimizeWaypoints optimizes the order of waypoints for the most efficient route
func (c *Client) OptimizeWaypoints(ctx context.Context, origin models.Location, waypoints []models.Location) ([]int, error) {
	if len(waypoints) == 0 {
		return []int{}, nil
	}

	waypointStrs := make([]string, len(waypoints))
	for i, wp := range waypoints {
		waypointStrs[i] = fmt.Sprintf("%f,%f", wp.Latitude, wp.Longitude)
	}

	req := &maps.DirectionsRequest{
		Origin:           fmt.Sprintf("%f,%f", origin.Latitude, origin.Longitude),
		Destination:      waypointStrs[len(waypointStrs)-1],
		Waypoints:        waypointStrs[:len(waypointStrs)-1],
		OptimizeWaypoints: true,
	}

	routes, _, err := c.client.Directions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("waypoint optimization failed: %w", err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no optimized route found")
	}

	return routes[0].WaypointOrder, nil
}

// CalculateDistanceMatrix calculates distances and durations between multiple origins and destinations
func (c *Client) CalculateDistanceMatrix(ctx context.Context, origins, destinations []models.Location) (*maps.DistanceMatrixResponse, error) {
	originStrs := make([]string, len(origins))
	for i, o := range origins {
		originStrs[i] = fmt.Sprintf("%f,%f", o.Latitude, o.Longitude)
	}

	destStrs := make([]string, len(destinations))
	for i, d := range destinations {
		destStrs[i] = fmt.Sprintf("%f,%f", d.Latitude, d.Longitude)
	}

	req := &maps.DistanceMatrixRequest{
		Origins:      originStrs,
		Destinations: destStrs,
	}

	resp, err := c.client.DistanceMatrix(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("distance matrix request failed: %w", err)
	}

	return resp, nil
}

// RouteInfo contains information about a calculated route
type RouteInfo struct {
	DistanceMeters  int64
	DurationSeconds int64
	Polyline        string
}
