package maps

import (
	"context"
	"fmt"

	"googlemaps.github.io/maps"
)

// Client represents the Google Maps API client
type Client struct {
	*maps.Client
}

// NewClient creates a new Google Maps client
func NewClient(apiKey string) (*Client, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %w", err)
	}

	return &Client{client}, nil
}

// Geocode converts an address to coordinates
func (c *Client) Geocode(address string) (*GeocodeResult, error) {
	req := &maps.GeocodingRequest{
		Address: address,
	}

	results, err := c.Client.Geocode(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	result := results[0]
	return &GeocodeResult{
		FormattedAddress: result.FormattedAddress,
		Latitude:         result.Geometry.Location.Lat,
		Longitude:        result.Geometry.Location.Lng,
		PlaceID:          result.PlaceID,
	}, nil
}

// ReverseGeocode converts coordinates to an address
func (c *Client) ReverseGeocode(lat, lon float64) (*GeocodeResult, error) {
	req := &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: lat,
			Lng: lon,
		},
	}

	results, err := c.Client.Geocode(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("reverse geocoding failed: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for coordinates: %f, %f", lat, lon)
	}

	result := results[0]
	return &GeocodeResult{
		FormattedAddress: result.FormattedAddress,
		Latitude:         result.Geometry.Location.Lat,
		Longitude:        result.Geometry.Location.Lng,
		PlaceID:          result.PlaceID,
	}, nil
}

// CalculateDistance calculates distance and duration between two points
func (c *Client) CalculateDistance(originLat, originLon, destLat, destLon float64) (*DistanceResult, error) {
	req := &maps.DistanceMatrixRequest{
		Origins: []string{
			fmt.Sprintf("%f,%f", originLat, originLon),
		},
		Destinations: []string{
			fmt.Sprintf("%f,%f", destLat, destLon),
		},
		Mode: maps.TravelModeDriving,
	}

	resp, err := c.Client.DistanceMatrix(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("distance matrix failed: %w", err)
	}

	if len(resp.Rows) == 0 || len(resp.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("no route found")
	}

	element := resp.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("route calculation failed: %s", element.Status)
	}

	return &DistanceResult{
		DistanceMeters: element.Distance.Meters,
		DistanceMiles:  float64(element.Distance.Meters) * 0.000621371,
		DurationSeconds: element.Duration.Seconds(),
		DurationText:   element.Duration.String(),
	}, nil
}

// GetDirections gets turn-by-turn directions between two points
func (c *Client) GetDirections(originLat, originLon, destLat, destLon float64) (*DirectionsResult, error) {
	req := &maps.DirectionsRequest{
		Origin: fmt.Sprintf("%f,%f", originLat, originLon),
		Destination: fmt.Sprintf("%f,%f", destLat, destLon),
		Mode: maps.TravelModeDriving,
	}

	routes, _, err := c.Client.Directions(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("directions failed: %w", err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	route := routes[0]
	var steps []DirectionStep

	for _, leg := range route.Legs {
		for _, step := range leg.Steps {
			steps = append(steps, DirectionStep{
				Instruction:     step.HTMLInstructions,
				DistanceMeters:  step.Distance.Meters,
				DurationSeconds: step.Duration.Seconds(),
				StartLat:        step.StartLocation.Lat,
				StartLon:        step.StartLocation.Lng,
				EndLat:          step.EndLocation.Lat,
				EndLon:          step.EndLocation.Lng,
			})
		}
	}

	totalDistance := 0
	totalDuration := 0.0
	for _, leg := range route.Legs {
		totalDistance += leg.Distance.Meters
		totalDuration += leg.Duration.Seconds()
	}

	return &DirectionsResult{
		DistanceMeters:  totalDistance,
		DistanceMiles:   float64(totalDistance) * 0.000621371,
		DurationSeconds: totalDuration,
		Polyline:        route.OverviewPolyline.Points,
		Steps:           steps,
	}, nil
}

// GeocodeResult represents geocoding results
type GeocodeResult struct {
	FormattedAddress string
	Latitude         float64
	Longitude        float64
	PlaceID          string
}

// DistanceResult represents distance calculation results
type DistanceResult struct {
	DistanceMeters  int
	DistanceMiles   float64
	DurationSeconds float64
	DurationText    string
}

// DirectionsResult represents turn-by-turn directions
type DirectionsResult struct {
	DistanceMeters  int
	DistanceMiles   float64
	DurationSeconds float64
	Polyline        string
	Steps           []DirectionStep
}

// DirectionStep represents a single navigation step
type DirectionStep struct {
	Instruction     string
	DistanceMeters  int
	DurationSeconds float64
	StartLat        float64
	StartLon        float64
	EndLat          float64
	EndLon          float64
}
