package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// DriverRepository handles driver database operations
type DriverRepository struct {
	db *DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *DB) *DriverRepository {
	return &DriverRepository{db: db}
}

// Create inserts a new driver
func (r *DriverRepository) Create(ctx context.Context, driver *models.Driver) error {
	query := `
		INSERT INTO drivers (id, name, phone_number, email, license_no, status, vehicle_id, motive_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		driver.ID, driver.Name, driver.PhoneNumber, driver.Email, driver.LicenseNo,
		driver.Status, driver.VehicleID, nil, driver.CreatedAt, driver.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}
	return nil
}

// GetByID retrieves a driver by ID
func (r *DriverRepository) GetByID(ctx context.Context, id string) (*models.Driver, error) {
	query := `
		SELECT id, name, phone_number, email, license_no, status, vehicle_id, created_at, updated_at
		FROM drivers WHERE id = ?
	`
	var driver models.Driver
	var vehicleID sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&driver.ID, &driver.Name, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo,
		&driver.Status, &vehicleID, &driver.CreatedAt, &driver.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("driver not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	if vehicleID.Valid {
		driver.VehicleID = vehicleID.String
	}

	return &driver, nil
}

// List retrieves all drivers
func (r *DriverRepository) List(ctx context.Context) ([]*models.Driver, error) {
	query := `
		SELECT id, name, phone_number, email, license_no, status, vehicle_id, created_at, updated_at
		FROM drivers ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*models.Driver
	for rows.Next() {
		var driver models.Driver
		var vehicleID sql.NullString

		if err := rows.Scan(
			&driver.ID, &driver.Name, &driver.PhoneNumber, &driver.Email, &driver.LicenseNo,
			&driver.Status, &vehicleID, &driver.CreatedAt, &driver.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}

		if vehicleID.Valid {
			driver.VehicleID = vehicleID.String
		}

		drivers = append(drivers, &driver)
	}

	return drivers, nil
}

// Update updates a driver
func (r *DriverRepository) Update(ctx context.Context, driver *models.Driver) error {
	query := `
		UPDATE drivers
		SET name = ?, phone_number = ?, email = ?, license_no = ?, status = ?, vehicle_id = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		driver.Name, driver.PhoneNumber, driver.Email, driver.LicenseNo,
		driver.Status, driver.VehicleID, driver.UpdatedAt, driver.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}
	return nil
}

// CreateFromMotive creates or updates a driver from Motive sync
func (r *DriverRepository) CreateFromMotive(ctx context.Context, driver *models.Driver, motiveID string) error {
	query := `
		INSERT INTO drivers (id, name, phone_number, email, license_no, status, motive_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(motive_id) DO UPDATE SET
			name = excluded.name,
			phone_number = excluded.phone_number,
			email = excluded.email,
			updated_at = excluded.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		driver.ID, driver.Name, driver.PhoneNumber, driver.Email, driver.LicenseNo,
		driver.Status, motiveID, driver.CreatedAt, driver.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create/update driver from Motive: %w", err)
	}
	return nil
}
