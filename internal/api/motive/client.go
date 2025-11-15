package motive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// Client wraps the Motive API client
type Client struct {
	baseURL   string
	apiKey    string
	apiSecret string
	client    *http.Client
}

// NewClient creates a new Motive API client
func NewClient(baseURL, apiKey, apiSecret string) *Client {
	return &Client{
		baseURL:   baseURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// MotiveVehicle represents a vehicle in the Motive system
type MotiveVehicle struct {
	ID           string  `json:"id"`
	Number       string  `json:"number"`
	Year         int     `json:"year"`
	Make         string  `json:"make"`
	Model        string  `json:"model"`
	VIN          string  `json:"vin"`
	LicensePlate string  `json:"license_plate"`
	CurrentLocation *MotiveLocation `json:"current_location,omitempty"`
}

// MotiveLocation represents a location in the Motive system
type MotiveLocation struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// MotiveDriver represents a driver in the Motive system
type MotiveDriver struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
}

// GetVehicles retrieves all vehicles from Motive
func (c *Client) GetVehicles(ctx context.Context) ([]MotiveVehicle, error) {
	req, err := c.newRequest(ctx, "GET", "/vehicles", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Vehicles []MotiveVehicle `json:"vehicles"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Vehicles, nil
}

// GetVehicle retrieves a specific vehicle from Motive
func (c *Client) GetVehicle(ctx context.Context, vehicleID string) (*MotiveVehicle, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/vehicles/%s", vehicleID), nil)
	if err != nil {
		return nil, err
	}

	var vehicle MotiveVehicle
	if err := c.doRequest(req, &vehicle); err != nil {
		return nil, err
	}

	return &vehicle, nil
}

// GetVehicleLocation retrieves the current location of a vehicle
func (c *Client) GetVehicleLocation(ctx context.Context, vehicleID string) (*models.Location, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/vehicles/%s/locations/latest", vehicleID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Location MotiveLocation `json:"location"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &models.Location{
		Latitude:  response.Location.Latitude,
		Longitude: response.Location.Longitude,
	}, nil
}

// GetDrivers retrieves all drivers from Motive
func (c *Client) GetDrivers(ctx context.Context) ([]MotiveDriver, error) {
	req, err := c.newRequest(ctx, "GET", "/drivers", nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Drivers []MotiveDriver `json:"drivers"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Drivers, nil
}

// GetDriver retrieves a specific driver from Motive
func (c *Client) GetDriver(ctx context.Context, driverID string) (*MotiveDriver, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/drivers/%s", driverID), nil)
	if err != nil {
		return nil, err
	}

	var driver MotiveDriver
	if err := c.doRequest(req, &driver); err != nil {
		return nil, err
	}

	return &driver, nil
}

// GetDriverLocation retrieves the current location of a driver
func (c *Client) GetDriverLocation(ctx context.Context, driverID string) (*models.Location, error) {
	req, err := c.newRequest(ctx, "GET", fmt.Sprintf("/drivers/%s/locations/latest", driverID), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Location MotiveLocation `json:"location"`
	}

	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &models.Location{
		Latitude:  response.Location.Latitude,
		Longitude: response.Location.Longitude,
	}, nil
}

// newRequest creates a new HTTP request with authentication
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("X-API-Secret", c.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// doRequest executes an HTTP request and decodes the response
func (c *Client) doRequest(req *http.Request, result interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
