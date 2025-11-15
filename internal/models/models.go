package models

import "time"

// Location represents a geographical point
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

// Driver represents a driver in the dispatch system
type Driver struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	LicenseNo   string    `json:"license_no"`
	Status      string    `json:"status"` // available, on_route, off_duty
	VehicleID   string    `json:"vehicle_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Vehicle represents a vehicle in the fleet
type Vehicle struct {
	ID           string    `json:"id"`
	VIN          string    `json:"vin"`
	LicensePlate string    `json:"license_plate"`
	Make         string    `json:"make"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	Status       string    `json:"status"` // active, maintenance, inactive
	Location     *Location `json:"location,omitempty"`
	DriverID     string    `json:"driver_id,omitempty"`
	MotiveID     string    `json:"motive_id,omitempty"` // ID in Motive system
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Stop represents a stop on a route
type Stop struct {
	ID              string     `json:"id"`
	Location        Location   `json:"location"`
	Type            string     `json:"type"` // pickup, delivery
	ScheduledTime   *time.Time `json:"scheduled_time,omitempty"`
	ArrivalTime     *time.Time `json:"arrival_time,omitempty"`
	DepartureTime   *time.Time `json:"departure_time,omitempty"`
	Status          string     `json:"status"` // pending, completed, skipped
	Notes           string     `json:"notes,omitempty"`
	ContactName     string     `json:"contact_name,omitempty"`
	ContactPhone    string     `json:"contact_phone,omitempty"`
	EstimatedDuration int      `json:"estimated_duration"` // in minutes
}

// Route represents a delivery/pickup route
type Route struct {
	ID              string     `json:"id"`
	DriverID        string     `json:"driver_id"`
	VehicleID       string     `json:"vehicle_id"`
	Stops           []Stop     `json:"stops"`
	Status          string     `json:"status"` // planned, in_progress, completed, cancelled
	StartTime       *time.Time `json:"start_time,omitempty"`
	EndTime         *time.Time `json:"end_time,omitempty"`
	TotalDistance   float64    `json:"total_distance"`   // in miles
	EstimatedDuration int      `json:"estimated_duration"` // in minutes
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// DispatchRequest represents a request to create a new dispatch route
type DispatchRequest struct {
	DriverID  string     `json:"driver_id"`
	VehicleID string     `json:"vehicle_id"`
	Stops     []Stop     `json:"stops"`
	StartTime *time.Time `json:"start_time,omitempty"`
}

// RouteOptimizationRequest represents a request to optimize a route
type RouteOptimizationRequest struct {
	Origin       Location   `json:"origin"`
	Destinations []Location `json:"destinations"`
	VehicleID    string     `json:"vehicle_id,omitempty"`
}

// RouteOptimizationResponse represents an optimized route
type RouteOptimizationResponse struct {
	OptimizedOrder []int     `json:"optimized_order"`
	TotalDistance  float64   `json:"total_distance"`
	TotalDuration  int       `json:"total_duration"`
	Route          *Route    `json:"route,omitempty"`
}
