package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// Connect establishes a connection to the PostgreSQL database
func Connect(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Vehicle represents a vehicle record in the database
type Vehicle struct {
	VehicleID        int
	VehicleNumber    string
	EntityType       *string
	CurrentLat       *float64
	CurrentLon       *float64
	SpeedMPH         *int
	Heading          *float64
	FuelGallons      *float64
	FuelPercent      *float64
	EngineHours      *int
	OdometerMiles    *int
	State            *string
	Status           *string
	VIN              *string
	Make             *string
	Model            *string
	Year             *int
	AssignedDriverID *int
	LastUpdated      time.Time
	CreatedAt        time.Time
}

// Asset represents an asset (trailer) record in the database
type Asset struct {
	AssetID             int
	AssetName           string
	AssetType           *string
	CurrentLat          *float64
	CurrentLon          *float64
	Status              *string
	LastMoved           *time.Time
	AttachedToVehicleID *int
	LocationDescription *string
	LastUpdated         time.Time
	CreatedAt           time.Time
}

// Driver represents a driver record in the database
type Driver struct {
	DriverID          int
	DriverCompanyID   *string
	FirstName         *string
	LastName          *string
	CurrentLat        *float64
	CurrentLon        *float64
	HOSRemainingMins  *int
	HOSStatus         *string
	CurrentVehicleID  *int
	HomeTerminalID    *int
	Phone             *string
	Email             *string
	LastUpdated       time.Time
	CreatedAt         time.Time
}

// Geofence represents a geofence record in the database
type Geofence struct {
	GeofenceID     int
	Name           string
	Description    *string
	BoundaryType   *string
	CenterLat      *float64
	CenterLon      *float64
	RadiusMeters   *int
	BoundaryPoints *string // JSON
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
