package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite database and creates tables
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS drivers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		email TEXT NOT NULL,
		license_number TEXT NOT NULL,
		status TEXT DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS vehicles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		make TEXT NOT NULL,
		model TEXT NOT NULL,
		year INTEGER NOT NULL,
		vin TEXT NOT NULL UNIQUE,
		license_plate TEXT NOT NULL,
		status TEXT DEFAULT 'available',
		current_lat REAL,
		current_lng REAL,
		last_updated DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS routes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		driver_id INTEGER NOT NULL,
		vehicle_id INTEGER NOT NULL,
		status TEXT DEFAULT 'planned',
		planned_start_time DATETIME NOT NULL,
		actual_start_time DATETIME,
		completed_time DATETIME,
		total_distance REAL DEFAULT 0,
		estimated_time INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (driver_id) REFERENCES drivers(id),
		FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
	);

	CREATE TABLE IF NOT EXISTS stops (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		route_id INTEGER NOT NULL,
		address TEXT NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		stop_type TEXT NOT NULL,
		stop_order INTEGER NOT NULL,
		status TEXT DEFAULT 'pending',
		notes TEXT,
		completed_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (route_id) REFERENCES routes(id)
	);

	CREATE INDEX IF NOT EXISTS idx_routes_driver ON routes(driver_id);
	CREATE INDEX IF NOT EXISTS idx_routes_vehicle ON routes(vehicle_id);
	CREATE INDEX IF NOT EXISTS idx_routes_status ON routes(status);
	CREATE INDEX IF NOT EXISTS idx_stops_route ON stops(route_id);
	`

	_, err := db.Exec(schema)
	return err
}
