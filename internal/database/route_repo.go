package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mcdougaj/Go-Dispatch/internal/models"
)

// RouteRepository handles route database operations
type RouteRepository struct {
	db *DB
}

// NewRouteRepository creates a new route repository
func NewRouteRepository(db *DB) *RouteRepository {
	return &RouteRepository{db: db}
}

// Create inserts a new route with its stops
func (r *RouteRepository) Create(ctx context.Context, route *models.Route) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert route
	routeQuery := `
		INSERT INTO routes (id, driver_id, vehicle_id, status, start_time, end_time,
			total_distance, estimated_duration, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, routeQuery,
		route.ID, route.DriverID, route.VehicleID, route.Status, route.StartTime, route.EndTime,
		route.TotalDistance, route.EstimatedDuration, route.CreatedAt, route.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create route: %w", err)
	}

	// Insert stops
	stopQuery := `
		INSERT INTO stops (id, route_id, stop_order, location_latitude, location_longitude, location_address,
			type, scheduled_time, arrival_time, departure_time, status, notes, contact_name, contact_phone,
			estimated_duration, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	for i, stop := range route.Stops {
		_, err = tx.ExecContext(ctx, stopQuery,
			stop.ID, route.ID, i, stop.Location.Latitude, stop.Location.Longitude, stop.Location.Address,
			stop.Type, stop.ScheduledTime, stop.ArrivalTime, stop.DepartureTime, stop.Status,
			stop.Notes, stop.ContactName, stop.ContactPhone, stop.EstimatedDuration,
		)
		if err != nil {
			return fmt.Errorf("failed to create stop: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a route with its stops
func (r *RouteRepository) GetByID(ctx context.Context, id string) (*models.Route, error) {
	routeQuery := `
		SELECT id, driver_id, vehicle_id, status, start_time, end_time,
			total_distance, estimated_duration, created_at, updated_at
		FROM routes WHERE id = ?
	`
	var route models.Route
	var startTime, endTime sql.NullTime

	err := r.db.QueryRowContext(ctx, routeQuery, id).Scan(
		&route.ID, &route.DriverID, &route.VehicleID, &route.Status, &startTime, &endTime,
		&route.TotalDistance, &route.EstimatedDuration, &route.CreatedAt, &route.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("route not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get route: %w", err)
	}

	if startTime.Valid {
		route.StartTime = &startTime.Time
	}
	if endTime.Valid {
		route.EndTime = &endTime.Time
	}

	// Get stops
	stopsQuery := `
		SELECT id, location_latitude, location_longitude, location_address, type,
			scheduled_time, arrival_time, departure_time, status, notes, contact_name,
			contact_phone, estimated_duration
		FROM stops WHERE route_id = ? ORDER BY stop_order
	`
	rows, err := r.db.QueryContext(ctx, stopsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get stops: %w", err)
	}
	defer rows.Close()

	var stops []models.Stop
	for rows.Next() {
		var stop models.Stop
		var addr, notes, contactName, contactPhone sql.NullString
		var scheduledTime, arrivalTime, departureTime sql.NullTime

		if err := rows.Scan(
			&stop.ID, &stop.Location.Latitude, &stop.Location.Longitude, &addr, &stop.Type,
			&scheduledTime, &arrivalTime, &departureTime, &stop.Status, &notes, &contactName,
			&contactPhone, &stop.EstimatedDuration,
		); err != nil {
			return nil, fmt.Errorf("failed to scan stop: %w", err)
		}

		if addr.Valid {
			stop.Location.Address = addr.String
		}
		if scheduledTime.Valid {
			stop.ScheduledTime = &scheduledTime.Time
		}
		if arrivalTime.Valid {
			stop.ArrivalTime = &arrivalTime.Time
		}
		if departureTime.Valid {
			stop.DepartureTime = &departureTime.Time
		}
		if notes.Valid {
			stop.Notes = notes.String
		}
		if contactName.Valid {
			stop.ContactName = contactName.String
		}
		if contactPhone.Valid {
			stop.ContactPhone = contactPhone.String
		}

		stops = append(stops, stop)
	}

	route.Stops = stops
	return &route, nil
}

// List retrieves all routes
func (r *RouteRepository) List(ctx context.Context) ([]*models.Route, error) {
	query := `
		SELECT id, driver_id, vehicle_id, status, start_time, end_time,
			total_distance, estimated_duration, created_at, updated_at
		FROM routes ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %w", err)
	}
	defer rows.Close()

	var routes []*models.Route
	for rows.Next() {
		var route models.Route
		var startTime, endTime sql.NullTime

		if err := rows.Scan(
			&route.ID, &route.DriverID, &route.VehicleID, &route.Status, &startTime, &endTime,
			&route.TotalDistance, &route.EstimatedDuration, &route.CreatedAt, &route.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan route: %w", err)
		}

		if startTime.Valid {
			route.StartTime = &startTime.Time
		}
		if endTime.Valid {
			route.EndTime = &endTime.Time
		}

		// Get stops for this route
		fullRoute, err := r.GetByID(ctx, route.ID)
		if err != nil {
			return nil, err
		}
		route.Stops = fullRoute.Stops

		routes = append(routes, &route)
	}

	return routes, nil
}

// UpdateStatus updates the status of a route
func (r *RouteRepository) UpdateStatus(ctx context.Context, id, status string, startTime, endTime *sql.NullTime) error {
	query := `
		UPDATE routes
		SET status = ?, start_time = ?, end_time = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, startTime, endTime, id)
	if err != nil {
		return fmt.Errorf("failed to update route status: %w", err)
	}
	return nil
}
