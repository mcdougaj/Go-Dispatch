package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// VehicleRepository handles vehicle database operations
type VehicleRepository struct {
	db *DB
}

// NewVehicleRepository creates a new vehicle repository
func NewVehicleRepository(db *DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// Create inserts a new vehicle
func (r *VehicleRepository) Create(ctx context.Context, vehicle *models.Vehicle) error {
	query := `
		INSERT INTO vehicles (id, vin, license_plate, make, model, year, status,
			location_latitude, location_longitude, location_address, driver_id, motive_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	var lat, lng sql.NullFloat64
	var addr sql.NullString
	if vehicle.Location != nil {
		lat = sql.NullFloat64{Float64: vehicle.Location.Latitude, Valid: true}
		lng = sql.NullFloat64{Float64: vehicle.Location.Longitude, Valid: true}
		addr = sql.NullString{String: vehicle.Location.Address, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		vehicle.ID, vehicle.VIN, vehicle.LicensePlate, vehicle.Make, vehicle.Model, vehicle.Year,
		vehicle.Status, lat, lng, addr, vehicle.DriverID, vehicle.MotiveID,
		vehicle.CreatedAt, vehicle.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}
	return nil
}

// GetByID retrieves a vehicle by ID
func (r *VehicleRepository) GetByID(ctx context.Context, id string) (*models.Vehicle, error) {
	query := `
		SELECT id, vin, license_plate, make, model, year, status,
			location_latitude, location_longitude, location_address, driver_id, motive_id, created_at, updated_at
		FROM vehicles WHERE id = ?
	`
	var vehicle models.Vehicle
	var lat, lng sql.NullFloat64
	var addr, driverID, motiveID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vehicle.ID, &vehicle.VIN, &vehicle.LicensePlate, &vehicle.Make, &vehicle.Model, &vehicle.Year,
		&vehicle.Status, &lat, &lng, &addr, &driverID, &motiveID,
		&vehicle.CreatedAt, &vehicle.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vehicle not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	if lat.Valid && lng.Valid {
		vehicle.Location = &models.Location{
			Latitude:  lat.Float64,
			Longitude: lng.Float64,
		}
		if addr.Valid {
			vehicle.Location.Address = addr.String
		}
	}

	if driverID.Valid {
		vehicle.DriverID = driverID.String
	}
	if motiveID.Valid {
		vehicle.MotiveID = motiveID.String
	}

	return &vehicle, nil
}

// List retrieves all vehicles
func (r *VehicleRepository) List(ctx context.Context) ([]*models.Vehicle, error) {
	query := `
		SELECT id, vin, license_plate, make, model, year, status,
			location_latitude, location_longitude, location_address, driver_id, motive_id, created_at, updated_at
		FROM vehicles ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		var lat, lng sql.NullFloat64
		var addr, driverID, motiveID sql.NullString

		if err := rows.Scan(
			&vehicle.ID, &vehicle.VIN, &vehicle.LicensePlate, &vehicle.Make, &vehicle.Model, &vehicle.Year,
			&vehicle.Status, &lat, &lng, &addr, &driverID, &motiveID,
			&vehicle.CreatedAt, &vehicle.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		if lat.Valid && lng.Valid {
			vehicle.Location = &models.Location{
				Latitude:  lat.Float64,
				Longitude: lng.Float64,
			}
			if addr.Valid {
				vehicle.Location.Address = addr.String
			}
		}

		if driverID.Valid {
			vehicle.DriverID = driverID.String
		}
		if motiveID.Valid {
			vehicle.MotiveID = motiveID.String
		}

		vehicles = append(vehicles, &vehicle)
	}

	return vehicles, nil
}

// Update updates a vehicle
func (r *VehicleRepository) Update(ctx context.Context, vehicle *models.Vehicle) error {
	query := `
		UPDATE vehicles
		SET vin = ?, license_plate = ?, make = ?, model = ?, year = ?, status = ?,
			location_latitude = ?, location_longitude = ?, location_address = ?, driver_id = ?, updated_at = ?
		WHERE id = ?
	`
	var lat, lng sql.NullFloat64
	var addr sql.NullString
	if vehicle.Location != nil {
		lat = sql.NullFloat64{Float64: vehicle.Location.Latitude, Valid: true}
		lng = sql.NullFloat64{Float64: vehicle.Location.Longitude, Valid: true}
		addr = sql.NullString{String: vehicle.Location.Address, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		vehicle.VIN, vehicle.LicensePlate, vehicle.Make, vehicle.Model, vehicle.Year,
		vehicle.Status, lat, lng, addr, vehicle.DriverID, vehicle.UpdatedAt, vehicle.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

// CreateFromMotive creates or updates a vehicle from Motive sync
func (r *VehicleRepository) CreateFromMotive(ctx context.Context, vehicle *models.Vehicle, motiveID string) error {
	query := `
		INSERT INTO vehicles (id, vin, license_plate, make, model, year, status,
			location_latitude, location_longitude, location_address, motive_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(motive_id) DO UPDATE SET
			vin = excluded.vin,
			license_plate = excluded.license_plate,
			make = excluded.make,
			model = excluded.model,
			year = excluded.year,
			location_latitude = excluded.location_latitude,
			location_longitude = excluded.location_longitude,
			location_address = excluded.location_address,
			updated_at = excluded.updated_at
	`
	var lat, lng sql.NullFloat64
	var addr sql.NullString
	if vehicle.Location != nil {
		lat = sql.NullFloat64{Float64: vehicle.Location.Latitude, Valid: true}
		lng = sql.NullFloat64{Float64: vehicle.Location.Longitude, Valid: true}
		addr = sql.NullString{String: vehicle.Location.Address, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		vehicle.ID, vehicle.VIN, vehicle.LicensePlate, vehicle.Make, vehicle.Model, vehicle.Year,
		vehicle.Status, lat, lng, addr, motiveID, vehicle.CreatedAt, vehicle.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create/update vehicle from Motive: %w", err)
	}
	return nil
}
