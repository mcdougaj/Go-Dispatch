package repositories

import (
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(route *models.CreateRouteRequest) (*models.Route, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert route
	query := `
		INSERT INTO routes (driver_id, vehicle_id, planned_start_time)
		VALUES (?, ?, ?)
	`
	result, err := tx.Exec(query, route.DriverID, route.VehicleID, route.PlannedStartTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create route: %w", err)
	}

	routeID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	// Insert stops
	stopQuery := `
		INSERT INTO stops (route_id, address, latitude, longitude, stop_type, stop_order, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	for i, stop := range route.Stops {
		_, err := tx.Exec(stopQuery, routeID, stop.Address, stop.Latitude, stop.Longitude, stop.StopType, i+1, stop.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to create stop: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetByID(int(routeID))
}

func (r *RouteRepository) GetByID(id int) (*models.Route, error) {
	query := `
		SELECT id, driver_id, vehicle_id, status, planned_start_time,
		       actual_start_time, completed_time, total_distance, estimated_time,
		       created_at, updated_at
		FROM routes
		WHERE id = ?
	`
	var route models.Route
	var actualStart, completed sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&route.ID, &route.DriverID, &route.VehicleID, &route.Status,
		&route.PlannedStartTime, &actualStart, &completed,
		&route.TotalDistance, &route.EstimatedTime,
		&route.CreatedAt, &route.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	if actualStart.Valid {
		route.ActualStartTime = actualStart.Time
	}
	if completed.Valid {
		route.CompletedTime = completed.Time
	}

	// Get stops
	stops, err := r.GetStopsByRouteID(id)
	if err != nil {
		return nil, err
	}
	route.Stops = stops

	return &route, nil
}

func (r *RouteRepository) GetAll() ([]models.Route, error) {
	query := `
		SELECT id, driver_id, vehicle_id, status, planned_start_time,
		       actual_start_time, completed_time, total_distance, estimated_time,
		       created_at, updated_at
		FROM routes
		ORDER BY planned_start_time DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query routes: %w", err)
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var route models.Route
		var actualStart, completed sql.NullTime

		err := rows.Scan(
			&route.ID, &route.DriverID, &route.VehicleID, &route.Status,
			&route.PlannedStartTime, &actualStart, &completed,
			&route.TotalDistance, &route.EstimatedTime,
			&route.CreatedAt, &route.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}

		if actualStart.Valid {
			route.ActualStartTime = actualStart.Time
		}
		if completed.Valid {
			route.CompletedTime = completed.Time
		}

		routes = append(routes, route)
	}
	return routes, nil
}

func (r *RouteRepository) GetStopsByRouteID(routeID int) ([]models.Stop, error) {
	query := `
		SELECT id, route_id, address, latitude, longitude, stop_type,
		       stop_order, status, notes, completed_at, created_at, updated_at
		FROM stops
		WHERE route_id = ?
		ORDER BY stop_order
	`
	rows, err := r.db.Query(query, routeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query stops: %w", err)
	}
	defer rows.Close()

	var stops []models.Stop
	for rows.Next() {
		var stop models.Stop
		var notes sql.NullString
		var completedAt sql.NullTime

		err := rows.Scan(
			&stop.ID, &stop.RouteID, &stop.Address, &stop.Latitude, &stop.Longitude,
			&stop.StopType, &stop.StopOrder, &stop.Status, &notes, &completedAt,
			&stop.CreatedAt, &stop.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stop: %w", err)
		}

		if notes.Valid {
			stop.Notes = notes.String
		}
		if completedAt.Valid {
			stop.CompletedAt = completedAt.Time
		}

		stops = append(stops, stop)
	}
	return stops, nil
}

func (r *RouteRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE routes SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *RouteRepository) UpdateRouteMetrics(id int, distance float64, duration int) error {
	query := `
		UPDATE routes
		SET total_distance = ?, estimated_time = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.Exec(query, distance, duration, id)
	return err
}
