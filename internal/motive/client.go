package motive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Motive API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Motive API client
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest performs a GET request to the Motive API
func (c *Client) makeRequest(endpoint string) ([]byte, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

// GetVehicleLocations retrieves current vehicle locations
func (c *Client) GetVehicleLocations() (*VehicleLocationsResponse, error) {
	data, err := c.makeRequest("/v3/vehicle_locations")
	if err != nil {
		return nil, err
	}

	var response VehicleLocationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse vehicle locations: %w", err)
	}

	return &response, nil
}

// GetVehicles retrieves vehicle details
func (c *Client) GetVehicles() (*VehiclesResponse, error) {
	data, err := c.makeRequest("/v1/vehicles")
	if err != nil {
		return nil, err
	}

	var response VehiclesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse vehicles: %w", err)
	}

	return &response, nil
}

// GetAssetLocations retrieves asset (trailer) locations
func (c *Client) GetAssetLocations() (*AssetLocationsResponse, error) {
	data, err := c.makeRequest("/v1/assets/locations")
	if err != nil {
		return nil, err
	}

	var response AssetLocationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse asset locations: %w", err)
	}

	return &response, nil
}

// GetDriverLocations retrieves driver locations
func (c *Client) GetDriverLocations() (*DriverLocationsResponse, error) {
	data, err := c.makeRequest("/v1/driver_locations")
	if err != nil {
		return nil, err
	}

	var response DriverLocationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse driver locations: %w", err)
	}

	return &response, nil
}

// GetDrivers retrieves driver roster
func (c *Client) GetDrivers() (*DriversResponse, error) {
	data, err := c.makeRequest("/v1/users")
	if err != nil {
		return nil, err
	}

	var response DriversResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse drivers: %w", err)
	}

	return &response, nil
}

// GetHOSAvailableTime retrieves driver HOS available time
func (c *Client) GetHOSAvailableTime() (*HOSResponse, error) {
	data, err := c.makeRequest("/v1/hos/available_time")
	if err != nil {
		return nil, err
	}

	var response HOSResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse HOS data: %w", err)
	}

	return &response, nil
}

// GetGeofences retrieves geofence definitions
func (c *Client) GetGeofences() (*GeofencesResponse, error) {
	data, err := c.makeRequest("/v1/geofences")
	if err != nil {
		return nil, err
	}

	var response GeofencesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse geofences: %w", err)
	}

	return &response, nil
}

// GetDispatches retrieves active dispatches
func (c *Client) GetDispatches() (*DispatchesResponse, error) {
	data, err := c.makeRequest("/v3/dispatches")
	if err != nil {
		return nil, err
	}

	var response DispatchesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse dispatches: %w", err)
	}

	return &response, nil
}

// GetUtilization retrieves vehicle utilization data
func (c *Client) GetUtilization() (*UtilizationResponse, error) {
	data, err := c.makeRequest("/v1/utilization")
	if err != nil {
		return nil, err
	}

	var response UtilizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse utilization: %w", err)
	}

	return &response, nil
}
