package motive

import "time"

// VehicleLocationsResponse represents the response from vehicle locations endpoint
type VehicleLocationsResponse struct {
	Vehicles []VehicleLocation `json:"vehicle_locations"`
	Pagination
}

// VehicleLocation represents a single vehicle's location data
type VehicleLocation struct {
	ID             int       `json:"id"`
	VehicleNumber  string    `json:"number"`
	EntityType     string    `json:"entity_type"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	Speed          int       `json:"speed"`          // mph
	Heading        float64   `json:"heading"`        // degrees
	FuelML         int       `json:"fuel_ml"`        // milliliters
	EngineHours    int       `json:"engine_hours"`
	Odometer       int       `json:"odometer"`       // meters
	State          string    `json:"state"`
	Status         string    `json:"status"`         // moving/idling/stopped
	DriverID       int       `json:"driver_id"`
	LastUpdated    time.Time `json:"last_updated"`
}

// FuelGallons converts fuel from milliliters to gallons
func (v *VehicleLocation) FuelGallons() float64 {
	return float64(v.FuelML) / 3785.41 // 1 gallon = 3785.41 mL
}

// OdometerMiles converts odometer from meters to miles
func (v *VehicleLocation) OdometerMiles() int {
	return int(float64(v.Odometer) * 0.000621371) // 1 meter = 0.000621371 miles
}

// VehiclesResponse represents the response from vehicles endpoint
type VehiclesResponse struct {
	Vehicles []Vehicle `json:"vehicles"`
	Pagination
}

// Vehicle represents detailed vehicle information
type Vehicle struct {
	ID            int    `json:"id"`
	Number        string `json:"number"`
	VIN           string `json:"vin"`
	Make          string `json:"make"`
	Model         string `json:"model"`
	Year          int    `json:"year"`
	FuelType      string `json:"fuel_type"`
	FuelCapacity  int    `json:"fuel_capacity_ml"`
	VehicleType   string `json:"vehicle_type"`
	Status        string `json:"status"`
}

// AssetLocationsResponse represents the response from asset locations endpoint
type AssetLocationsResponse struct {
	Assets []AssetLocation `json:"assets"`
	Pagination
}

// AssetLocation represents a single asset's (trailer/container) location
type AssetLocation struct {
	ID                  int       `json:"id"`
	Name                string    `json:"name"`
	AssetType           string    `json:"asset_type"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	Status              string    `json:"status"` // loaded/empty
	LastMoved           time.Time `json:"last_moved"`
	AttachedToVehicleID *int      `json:"attached_to_vehicle_id"`
	Location            string    `json:"location_description"`
}

// DriverLocationsResponse represents the response from driver locations endpoint
type DriverLocationsResponse struct {
	Drivers []DriverLocation `json:"driver_locations"`
	Pagination
}

// DriverLocation represents a single driver's location
type DriverLocation struct {
	DriverID         int       `json:"driver_id"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	VehicleID        *int      `json:"vehicle_id"`
	LastUpdated      time.Time `json:"last_updated"`
}

// DriversResponse represents the response from users (drivers) endpoint
type DriversResponse struct {
	Users []Driver `json:"users"`
	Pagination
}

// Driver represents a driver's details
type Driver struct {
	ID              int    `json:"id"`
	CompanyID       string `json:"company_id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Status          string `json:"status"`
	HomeTerminalID  *int   `json:"home_terminal_id"`
	CurrentVehicleID *int  `json:"current_vehicle_id"`
}

// HOSResponse represents the response from HOS available time endpoint
type HOSResponse struct {
	HOSData []HOSAvailableTime `json:"hos_available_time"`
	Pagination
}

// HOSAvailableTime represents a driver's hours of service data
type HOSAvailableTime struct {
	DriverID              int    `json:"driver_id"`
	DriveTimeRemaining    int    `json:"drive_time_remaining_minutes"`
	ShiftTimeRemaining    int    `json:"shift_time_remaining_minutes"`
	CycleTimeRemaining    int    `json:"cycle_time_remaining_minutes"`
	Status                string `json:"status"` // driving/on_duty/sleeper/off_duty
	BreakRequired         bool   `json:"break_required"`
}

// GeofencesResponse represents the response from geofences endpoint
type GeofencesResponse struct {
	Geofences []Geofence `json:"geofences"`
	Pagination
}

// Geofence represents a geographic boundary
type Geofence struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"` // circle/polygon
	CenterLat   float64           `json:"center_latitude"`
	CenterLon   float64           `json:"center_longitude"`
	Radius      int               `json:"radius_meters"`
	Boundary    []BoundaryPoint   `json:"boundary_points"`
}

// BoundaryPoint represents a point in a polygon geofence
type BoundaryPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DispatchesResponse represents the response from dispatches endpoint
type DispatchesResponse struct {
	Dispatches []Dispatch `json:"dispatches"`
	Pagination
}

// Dispatch represents an active dispatch
type Dispatch struct {
	ID                int       `json:"id"`
	DriverID          int       `json:"driver_id"`
	VehicleID         int       `json:"vehicle_id"`
	Status            string    `json:"status"`
	OriginLat         float64   `json:"origin_latitude"`
	OriginLon         float64   `json:"origin_longitude"`
	OriginDescription string    `json:"origin_description"`
	DestLat           float64   `json:"destination_latitude"`
	DestLon           float64   `json:"destination_longitude"`
	DestDescription   string    `json:"destination_description"`
	ScheduledStart    time.Time `json:"scheduled_start"`
	ScheduledEnd      time.Time `json:"scheduled_end"`
}

// UtilizationResponse represents the response from utilization endpoint
type UtilizationResponse struct {
	Utilization []VehicleUtilization `json:"utilization"`
	Pagination
}

// VehicleUtilization represents vehicle usage metrics
type VehicleUtilization struct {
	VehicleID     int     `json:"vehicle_id"`
	DrivingTime   int     `json:"driving_time_minutes"`
	IdlingTime    int     `json:"idling_time_minutes"`
	StoppedTime   int     `json:"stopped_time_minutes"`
	Distance      int     `json:"distance_meters"`
	FuelConsumed  int     `json:"fuel_consumed_ml"`
}

// DistanceMiles converts distance from meters to miles
func (u *VehicleUtilization) DistanceMiles() float64 {
	return float64(u.Distance) * 0.000621371
}

// FuelConsumedGallons converts fuel from milliliters to gallons
func (u *VehicleUtilization) FuelConsumedGallons() float64 {
	return float64(u.FuelConsumed) / 3785.41
}

// Pagination represents common pagination fields in Motive API responses
type Pagination struct {
	PerPage     int `json:"per_page"`
	PageNo      int `json:"page_no"`
	TotalPages  int `json:"total_pages"`
	TotalRecords int `json:"total"`
}
