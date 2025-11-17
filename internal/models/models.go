package models

import "time"

// Driver represents a delivery driver
type Driver struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	LicenseNum  string    `json:"license_number"`
	Status      string    `json:"status"` // active, inactive, on_route
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Vehicle represents a delivery vehicle
type Vehicle struct {
	ID            int       `json:"id"`
	Make          string    `json:"make"`
	Model         string    `json:"model"`
	Year          int       `json:"year"`
	VIN           string    `json:"vin"`
	LicensePlate  string    `json:"license_plate"`
	Status        string    `json:"status"` // available, in_use, maintenance
	CurrentLat    float64   `json:"current_lat,omitempty"`
	CurrentLng    float64   `json:"current_lng,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Stop represents a delivery stop
type Stop struct {
	ID          int       `json:"id"`
	RouteID     int       `json:"route_id"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	StopType    string    `json:"stop_type"` // pickup, delivery
	StopOrder   int       `json:"stop_order"`
	Status      string    `json:"status"` // pending, completed, failed
	Notes       string    `json:"notes,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Route represents a delivery route
type Route struct {
	ID               int       `json:"id"`
	DriverID         int       `json:"driver_id"`
	VehicleID        int       `json:"vehicle_id"`
	Status           string    `json:"status"` // planned, in_progress, completed, cancelled
	PlannedStartTime time.Time `json:"planned_start_time"`
	ActualStartTime  time.Time `json:"actual_start_time,omitempty"`
	CompletedTime    time.Time `json:"completed_time,omitempty"`
	TotalDistance    float64   `json:"total_distance"` // in kilometers
	EstimatedTime    int       `json:"estimated_time"` // in minutes
	Stops            []Stop    `json:"stops,omitempty"`
	Driver           *Driver   `json:"driver,omitempty"`
	Vehicle          *Vehicle  `json:"vehicle,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateDriverRequest represents the request to create a driver
type CreateDriverRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	LicenseNum  string `json:"license_number"`
}

// CreateVehicleRequest represents the request to create a vehicle
type CreateVehicleRequest struct {
	Make         string `json:"make"`
	Model        string `json:"model"`
	Year         int    `json:"year"`
	VIN          string `json:"vin"`
	LicensePlate string `json:"license_plate"`
}

// CreateStopRequest represents a stop in a route creation request
type CreateStopRequest struct {
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	StopType  string  `json:"stop_type"`
	Notes     string  `json:"notes,omitempty"`
}

// CreateRouteRequest represents the request to create a route
type CreateRouteRequest struct {
	DriverID         int                 `json:"driver_id"`
	VehicleID        int                 `json:"vehicle_id"`
	PlannedStartTime time.Time           `json:"planned_start_time"`
	Stops            []CreateStopRequest `json:"stops"`
}

// UpdateRouteStatusRequest represents the request to update route status
type UpdateRouteStatusRequest struct {
	Status string `json:"status"`
}

// VehicleLocation represents real-time vehicle location
type VehicleLocation struct {
	VehicleID   int       `json:"vehicle_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Speed       float64   `json:"speed,omitempty"`
	Heading     float64   `json:"heading,omitempty"`
	LastUpdated time.Time `json:"last_updated"`
}
