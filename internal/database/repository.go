package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mcdougaj/go-dispatch/internal/motive"
)

// UpsertVehicleLocation inserts or updates a vehicle location
func (db *DB) UpsertVehicleLocation(v *motive.VehicleLocation) error {
	query := `
		INSERT INTO vehicles (
			vehicle_id, vehicle_number, entity_type, current_lat, current_lon,
			speed_mph, heading, fuel_gallons, engine_hours, odometer_miles,
			state, status, assigned_driver_id, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (vehicle_id) DO UPDATE SET
			vehicle_number = EXCLUDED.vehicle_number,
			entity_type = EXCLUDED.entity_type,
			current_lat = EXCLUDED.current_lat,
			current_lon = EXCLUDED.current_lon,
			speed_mph = EXCLUDED.speed_mph,
			heading = EXCLUDED.heading,
			fuel_gallons = EXCLUDED.fuel_gallons,
			engine_hours = EXCLUDED.engine_hours,
			odometer_miles = EXCLUDED.odometer_miles,
			state = EXCLUDED.state,
			status = EXCLUDED.status,
			assigned_driver_id = EXCLUDED.assigned_driver_id,
			last_updated = EXCLUDED.last_updated
	`

	var driverID *int
	if v.DriverID != 0 {
		driverID = &v.DriverID
	}

	_, err := db.Exec(query,
		v.ID,
		v.VehicleNumber,
		v.EntityType,
		v.Latitude,
		v.Longitude,
		v.Speed,
		v.Heading,
		v.FuelGallons(),
		v.EngineHours,
		v.OdometerMiles(),
		v.State,
		v.Status,
		driverID,
		v.LastUpdated,
	)

	return err
}

// UpsertVehicle inserts or updates vehicle details
func (db *DB) UpsertVehicle(v *motive.Vehicle) error {
	query := `
		INSERT INTO vehicles (
			vehicle_id, vehicle_number, vin, make, model, year
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (vehicle_id) DO UPDATE SET
			vehicle_number = EXCLUDED.vehicle_number,
			vin = EXCLUDED.vin,
			make = EXCLUDED.make,
			model = EXCLUDED.model,
			year = EXCLUDED.year
	`

	_, err := db.Exec(query, v.ID, v.Number, v.VIN, v.Make, v.Model, v.Year)
	return err
}

// UpsertAssetLocation inserts or updates an asset location
func (db *DB) UpsertAssetLocation(a *motive.AssetLocation) error {
	query := `
		INSERT INTO assets (
			asset_id, asset_name, asset_type, current_lat, current_lon,
			status, last_moved, attached_to_vehicle_id, location_description, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (asset_id) DO UPDATE SET
			asset_name = EXCLUDED.asset_name,
			asset_type = EXCLUDED.asset_type,
			current_lat = EXCLUDED.current_lat,
			current_lon = EXCLUDED.current_lon,
			status = EXCLUDED.status,
			last_moved = EXCLUDED.last_moved,
			attached_to_vehicle_id = EXCLUDED.attached_to_vehicle_id,
			location_description = EXCLUDED.location_description,
			last_updated = EXCLUDED.last_updated
	`

	_, err := db.Exec(query,
		a.ID,
		a.Name,
		a.AssetType,
		a.Latitude,
		a.Longitude,
		a.Status,
		a.LastMoved,
		a.AttachedToVehicleID,
		a.Location,
		time.Now(),
	)

	return err
}

// UpsertDriver inserts or updates driver details
func (db *DB) UpsertDriver(d *motive.Driver) error {
	query := `
		INSERT INTO drivers (
			driver_id, driver_company_id, first_name, last_name,
			phone, email, current_vehicle_id, home_terminal_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (driver_id) DO UPDATE SET
			driver_company_id = EXCLUDED.driver_company_id,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			current_vehicle_id = EXCLUDED.current_vehicle_id,
			home_terminal_id = EXCLUDED.home_terminal_id
	`

	_, err := db.Exec(query,
		d.ID,
		d.CompanyID,
		d.FirstName,
		d.LastName,
		d.Phone,
		d.Email,
		d.CurrentVehicleID,
		d.HomeTerminalID,
	)

	return err
}

// UpsertDriverLocation inserts or updates a driver location
func (db *DB) UpsertDriverLocation(dl *motive.DriverLocation) error {
	query := `
		UPDATE drivers SET
			current_lat = $1,
			current_lon = $2,
			current_vehicle_id = $3,
			last_updated = $4
		WHERE driver_id = $5
	`

	_, err := db.Exec(query,
		dl.Latitude,
		dl.Longitude,
		dl.VehicleID,
		dl.LastUpdated,
		dl.DriverID,
	)

	return err
}

// UpsertHOSData inserts or updates HOS data for a driver
func (db *DB) UpsertHOSData(hos *motive.HOSAvailableTime) error {
	query := `
		UPDATE drivers SET
			hos_remaining_minutes = $1,
			hos_status = $2,
			last_updated = $3
		WHERE driver_id = $4
	`

	_, err := db.Exec(query,
		hos.DriveTimeRemaining,
		hos.Status,
		time.Now(),
		hos.DriverID,
	)

	return err
}

// UpsertGeofence inserts or updates a geofence
func (db *DB) UpsertGeofence(g *motive.Geofence) error {
	var boundaryJSON []byte
	var err error
	if len(g.Boundary) > 0 {
		boundaryJSON, err = json.Marshal(g.Boundary)
		if err != nil {
			return fmt.Errorf("failed to marshal boundary points: %w", err)
		}
	}

	query := `
		INSERT INTO geofences (
			geofence_id, name, description, boundary_type,
			center_lat, center_lon, radius_meters, boundary_points, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (geofence_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			boundary_type = EXCLUDED.boundary_type,
			center_lat = EXCLUDED.center_lat,
			center_lon = EXCLUDED.center_lon,
			radius_meters = EXCLUDED.radius_meters,
			boundary_points = EXCLUDED.boundary_points,
			updated_at = EXCLUDED.updated_at
	`

	_, err = db.Exec(query,
		g.ID,
		g.Name,
		g.Description,
		g.Type,
		g.CenterLat,
		g.CenterLon,
		g.Radius,
		boundaryJSON,
		time.Now(),
	)

	return err
}

// GetAllVehicles retrieves all vehicles
func (db *DB) GetAllVehicles() ([]Vehicle, error) {
	query := `
		SELECT vehicle_id, vehicle_number, entity_type, current_lat, current_lon,
			speed_mph, heading, fuel_gallons, fuel_percent, engine_hours, odometer_miles,
			state, status, vin, make, model, year, assigned_driver_id, last_updated, created_at
		FROM vehicles
		ORDER BY vehicle_number
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []Vehicle
	for rows.Next() {
		var v Vehicle
		err := rows.Scan(
			&v.VehicleID, &v.VehicleNumber, &v.EntityType, &v.CurrentLat, &v.CurrentLon,
			&v.SpeedMPH, &v.Heading, &v.FuelGallons, &v.FuelPercent, &v.EngineHours, &v.OdometerMiles,
			&v.State, &v.Status, &v.VIN, &v.Make, &v.Model, &v.Year, &v.AssignedDriverID,
			&v.LastUpdated, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, rows.Err()
}

// GetAllAssets retrieves all assets
func (db *DB) GetAllAssets() ([]Asset, error) {
	query := `
		SELECT asset_id, asset_name, asset_type, current_lat, current_lon,
			status, last_moved, attached_to_vehicle_id, location_description,
			last_updated, created_at
		FROM assets
		ORDER BY asset_name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		err := rows.Scan(
			&a.AssetID, &a.AssetName, &a.AssetType, &a.CurrentLat, &a.CurrentLon,
			&a.Status, &a.LastMoved, &a.AttachedToVehicleID, &a.LocationDescription,
			&a.LastUpdated, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}

	return assets, rows.Err()
}

// GetAllDrivers retrieves all drivers
func (db *DB) GetAllDrivers() ([]Driver, error) {
	query := `
		SELECT driver_id, driver_company_id, first_name, last_name,
			current_lat, current_lon, hos_remaining_minutes, hos_status,
			current_vehicle_id, home_terminal_id, phone, email,
			last_updated, created_at
		FROM drivers
		ORDER BY driver_id
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var d Driver
		err := rows.Scan(
			&d.DriverID, &d.DriverCompanyID, &d.FirstName, &d.LastName,
			&d.CurrentLat, &d.CurrentLon, &d.HOSRemainingMins, &d.HOSStatus,
			&d.CurrentVehicleID, &d.HomeTerminalID, &d.Phone, &d.Email,
			&d.LastUpdated, &d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}

	return drivers, rows.Err()
}

// GetAllGeofences retrieves all geofences
func (db *DB) GetAllGeofences() ([]Geofence, error) {
	query := `
		SELECT geofence_id, name, description, boundary_type,
			center_lat, center_lon, radius_meters, boundary_points,
			created_at, updated_at
		FROM geofences
		ORDER BY name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var geofences []Geofence
	for rows.Next() {
		var g Geofence
		err := rows.Scan(
			&g.GeofenceID, &g.Name, &g.Description, &g.BoundaryType,
			&g.CenterLat, &g.CenterLon, &g.RadiusMeters, &g.BoundaryPoints,
			&g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		geofences = append(geofences, g)
	}

	return geofences, rows.Err()
}

// RecordSyncStatus records the status of a data sync operation
func (db *DB) RecordSyncStatus(entityType string, recordsUpdated int, status string, errorMsg *string) error {
	query := `
		INSERT INTO sync_status (entity_type, last_sync_time, records_updated, status, error_message)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(query, entityType, time.Now(), recordsUpdated, status, errorMsg)
	return err
}

// LogQuery logs a natural language query and AI response
func (db *DB) LogQuery(userQuery, aiResponse, queryIntent string, entitiesJSON []byte, responseTimeMS int) error {
	query := `
		INSERT INTO query_log (user_query, ai_response, query_intent, entities_extracted, response_time_ms)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(query, userQuery, aiResponse, queryIntent, entitiesJSON, responseTimeMS)
	return err
}
