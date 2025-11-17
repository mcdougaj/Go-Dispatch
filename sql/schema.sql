-- Motive Dispatch Application Database Schema
-- Read-only tracking and viewing application

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Vehicle tracking table
CREATE TABLE IF NOT EXISTS vehicles (
    vehicle_id INTEGER PRIMARY KEY,
    vehicle_number VARCHAR(50) NOT NULL,
    entity_type VARCHAR(20),
    current_lat DECIMAL(10, 8),
    current_lon DECIMAL(11, 8),
    speed_mph INTEGER,
    heading DECIMAL(5, 2),
    fuel_gallons DECIMAL(8, 2),  -- Converted from mL
    fuel_percent DECIMAL(5, 2),
    engine_hours INTEGER,
    odometer_miles INTEGER,
    state VARCHAR(50),
    status VARCHAR(20),  -- moving/idling/stopped
    vin VARCHAR(17),
    make VARCHAR(50),
    model VARCHAR(50),
    year INTEGER,
    assigned_driver_id INTEGER,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Asset (trailer/container) tracking table
CREATE TABLE IF NOT EXISTS assets (
    asset_id INTEGER PRIMARY KEY,
    asset_name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(50),  -- trailer, container, etc.
    current_lat DECIMAL(10, 8),
    current_lon DECIMAL(11, 8),
    status VARCHAR(20),  -- loaded/empty/maintenance
    last_moved TIMESTAMP,
    attached_to_vehicle_id INTEGER,
    location_description TEXT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (attached_to_vehicle_id) REFERENCES vehicles(vehicle_id)
);

-- Driver tracking table
CREATE TABLE IF NOT EXISTS drivers (
    driver_id INTEGER PRIMARY KEY,
    driver_company_id VARCHAR(50),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    current_lat DECIMAL(10, 8),
    current_lon DECIMAL(11, 8),
    hos_remaining_minutes INTEGER,
    hos_status VARCHAR(20),  -- available/driving/sleeper/off_duty
    current_vehicle_id INTEGER,
    home_terminal_id INTEGER,
    phone VARCHAR(20),
    email VARCHAR(100),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (current_vehicle_id) REFERENCES vehicles(vehicle_id)
);

-- Facilities table (yards, customer sites)
CREATE TABLE IF NOT EXISTS facilities (
    facility_id SERIAL PRIMARY KEY,
    facility_name VARCHAR(100) NOT NULL,
    facility_type VARCHAR(50),  -- yard/customer/vendor
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(2),
    zip VARCHAR(10),
    lat DECIMAL(10, 8),
    lon DECIMAL(11, 8),
    geofence_id INTEGER,  -- Link to Motive geofence
    operating_hours JSONB,
    contact_info JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Geofences table (from Motive)
CREATE TABLE IF NOT EXISTS geofences (
    geofence_id INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    boundary_type VARCHAR(20),  -- circle/polygon
    center_lat DECIMAL(10, 8),
    center_lon DECIMAL(11, 8),
    radius_meters INTEGER,
    boundary_points JSONB,  -- For polygon geofences
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Query log table (for AI interactions)
CREATE TABLE IF NOT EXISTS query_log (
    query_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_query TEXT NOT NULL,
    ai_response TEXT,
    query_intent VARCHAR(50),
    entities_extracted JSONB,
    response_time_ms INTEGER,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sync status table (track data synchronization)
CREATE TABLE IF NOT EXISTS sync_status (
    sync_id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,  -- vehicles/drivers/assets/geofences
    last_sync_time TIMESTAMP,
    records_updated INTEGER,
    status VARCHAR(20),  -- success/failed/in_progress
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_vehicles_location ON vehicles(current_lat, current_lon);
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_vehicles_driver ON vehicles(assigned_driver_id);
CREATE INDEX idx_vehicles_updated ON vehicles(last_updated);

CREATE INDEX idx_assets_location ON assets(current_lat, current_lon);
CREATE INDEX idx_assets_vehicle ON assets(attached_to_vehicle_id);
CREATE INDEX idx_assets_status ON assets(status);

CREATE INDEX idx_drivers_location ON drivers(current_lat, current_lon);
CREATE INDEX idx_drivers_vehicle ON drivers(current_vehicle_id);
CREATE INDEX idx_drivers_hos ON drivers(hos_status);

CREATE INDEX idx_facilities_location ON facilities(lat, lon);
CREATE INDEX idx_facilities_type ON facilities(facility_type);

CREATE INDEX idx_query_log_timestamp ON query_log(timestamp);
CREATE INDEX idx_sync_status_entity ON sync_status(entity_type, last_sync_time);

-- Comments for documentation
COMMENT ON TABLE vehicles IS 'Real-time vehicle location and status data from Motive API';
COMMENT ON TABLE assets IS 'Trailer and container tracking data';
COMMENT ON TABLE drivers IS 'Driver location and HOS (Hours of Service) data';
COMMENT ON TABLE facilities IS 'Company facilities and customer locations';
COMMENT ON TABLE geofences IS 'Geographic boundaries for facilities';
COMMENT ON TABLE query_log IS 'Natural language queries and AI responses';
COMMENT ON TABLE sync_status IS 'Data synchronization tracking from Motive API';
