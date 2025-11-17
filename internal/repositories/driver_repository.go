package repositories

import (
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository(db *sql.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) Create(driver *models.CreateDriverRequest) (*models.Driver, error) {
	query := `
		INSERT INTO drivers (name, phone_number, email, license_number)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, driver.Name, driver.PhoneNumber, driver.Email, driver.LicenseNum)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(int(id))
}

func (r *DriverRepository) GetByID(id int) (*models.Driver, error) {
	query := `
		SELECT id, name, phone_number, email, license_number, status, created_at, updated_at
		FROM drivers
		WHERE id = ?
	`
	var driver models.Driver
	err := r.db.QueryRow(query, id).Scan(
		&driver.ID, &driver.Name, &driver.PhoneNumber, &driver.Email,
		&driver.LicenseNum, &driver.Status, &driver.CreatedAt, &driver.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}
	return &driver, nil
}

func (r *DriverRepository) GetAll() ([]models.Driver, error) {
	query := `
		SELECT id, name, phone_number, email, license_number, status, created_at, updated_at
		FROM drivers
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var driver models.Driver
		err := rows.Scan(
			&driver.ID, &driver.Name, &driver.PhoneNumber, &driver.Email,
			&driver.LicenseNum, &driver.Status, &driver.CreatedAt, &driver.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}
	return drivers, nil
}

func (r *DriverRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE drivers SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}
