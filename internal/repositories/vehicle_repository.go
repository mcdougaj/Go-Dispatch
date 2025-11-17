package repositories

import (
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(vehicle *models.CreateVehicleRequest) (*models.Vehicle, error) {
	query := `
		INSERT INTO vehicles (make, model, year, vin, license_plate)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.VIN, vehicle.LicensePlate)
	if err != nil {
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(int(id))
}

func (r *VehicleRepository) GetByID(id int) (*models.Vehicle, error) {
	query := `
		SELECT id, make, model, year, vin, license_plate, status,
		       current_lat, current_lng, last_updated, created_at, updated_at
		FROM vehicles
		WHERE id = ?
	`
	var vehicle models.Vehicle
	var lat, lng sql.NullFloat64
	var lastUpdated sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year,
		&vehicle.VIN, &vehicle.LicensePlate, &vehicle.Status,
		&lat, &lng, &lastUpdated, &vehicle.CreatedAt, &vehicle.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle: %w", err)
	}

	if lat.Valid {
		vehicle.CurrentLat = lat.Float64
	}
	if lng.Valid {
		vehicle.CurrentLng = lng.Float64
	}
	if lastUpdated.Valid {
		vehicle.LastUpdated = lastUpdated.Time
	}

	return &vehicle, nil
}

func (r *VehicleRepository) GetAll() ([]models.Vehicle, error) {
	query := `
		SELECT id, make, model, year, vin, license_plate, status,
		       current_lat, current_lng, last_updated, created_at, updated_at
		FROM vehicles
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var vehicle models.Vehicle
		var lat, lng sql.NullFloat64
		var lastUpdated sql.NullTime

		err := rows.Scan(
			&vehicle.ID, &vehicle.Make, &vehicle.Model, &vehicle.Year,
			&vehicle.VIN, &vehicle.LicensePlate, &vehicle.Status,
			&lat, &lng, &lastUpdated, &vehicle.CreatedAt, &vehicle.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		if lat.Valid {
			vehicle.CurrentLat = lat.Float64
		}
		if lng.Valid {
			vehicle.CurrentLng = lng.Float64
		}
		if lastUpdated.Valid {
			vehicle.LastUpdated = lastUpdated.Time
		}

		vehicles = append(vehicles, vehicle)
	}
	return vehicles, nil
}

func (r *VehicleRepository) UpdateLocation(id int, lat, lng float64) error {
	query := `
		UPDATE vehicles
		SET current_lat = ?, current_lng = ?, last_updated = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, lat, lng, id)
	return err
}

func (r *VehicleRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE vehicles SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}
