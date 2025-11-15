-- Drivers table
CREATE TABLE IF NOT EXISTS drivers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT,
    email TEXT,
    license_no TEXT,
    status TEXT NOT NULL DEFAULT 'available',
    vehicle_id TEXT,
    motive_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
    id TEXT PRIMARY KEY,
    vin TEXT,
    license_plate TEXT,
    make TEXT,
    model TEXT,
    year INTEGER,
    status TEXT NOT NULL DEFAULT 'active',
    location_latitude REAL,
    location_longitude REAL,
    location_address TEXT,
    driver_id TEXT,
    motive_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (driver_id) REFERENCES drivers(id)
);

-- Routes table
CREATE TABLE IF NOT EXISTS routes (
    id TEXT PRIMARY KEY,
    driver_id TEXT NOT NULL,
    vehicle_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'planned',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    total_distance REAL DEFAULT 0,
    estimated_duration INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (driver_id) REFERENCES drivers(id),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
);

-- Stops table
CREATE TABLE IF NOT EXISTS stops (
    id TEXT PRIMARY KEY,
    route_id TEXT NOT NULL,
    stop_order INTEGER NOT NULL,
    location_latitude REAL NOT NULL,
    location_longitude REAL NOT NULL,
    location_address TEXT,
    type TEXT NOT NULL,
    scheduled_time TIMESTAMP,
    arrival_time TIMESTAMP,
    departure_time TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'pending',
    notes TEXT,
    contact_name TEXT,
    contact_phone TEXT,
    estimated_duration INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (route_id) REFERENCES routes(id) ON DELETE CASCADE
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_drivers_status ON drivers(status);
CREATE INDEX IF NOT EXISTS idx_drivers_motive_id ON drivers(motive_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_status ON vehicles(status);
CREATE INDEX IF NOT EXISTS idx_vehicles_motive_id ON vehicles(motive_id);
CREATE INDEX IF NOT EXISTS idx_routes_driver ON routes(driver_id);
CREATE INDEX IF NOT EXISTS idx_routes_vehicle ON routes(vehicle_id);
CREATE INDEX IF NOT EXISTS idx_routes_status ON routes(status);
CREATE INDEX IF NOT EXISTS idx_stops_route ON stops(route_id);
