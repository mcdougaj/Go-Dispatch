# Go-Dispatch

A comprehensive logistics and fleet management system built with Go and modern web technologies.

## Features

- **Driver Management**: Add, track, and manage delivery drivers with contact information and license details
- **Vehicle Fleet Management**: Monitor and manage your vehicle fleet with real-time status tracking
- **Route Planning & Optimization**: Create optimized delivery routes using Google Maps API
- **Real-time Tracking**: Track vehicle locations and route progress
- **Dispatch Dashboard**: Comprehensive dashboard for monitoring all operations
- **RESTful API**: Full-featured API for integration with other systems

## Technology Stack

- **Backend**: Go 1.21+
- **Database**: SQLite (pure Go, no CGO required)
- **Web Framework**: Gorilla Mux
- **Maps Integration**: Google Maps API
- **Frontend**: HTML5, CSS3, Vanilla JavaScript

## Prerequisites

- Go 1.21 or higher
- Google Maps API Key (required)
- Motive API credentials (optional)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/mcdougaj/Go-Dispatch.git
cd Go-Dispatch
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file from the example:
```bash
cp .env.example .env
```

4. Edit `.env` and add your API keys:
```env
SERVER_PORT=8080
GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here
MOTIVE_API_KEY=your_motive_api_key_here (optional)
MOTIVE_API_SECRET=your_motive_api_secret_here (optional)
DATABASE_PATH=./dispatch.db
```

## Running the Application

```bash
go run main.go
```

The application will be available at:
- Web Interface: http://localhost:8080
- API Endpoints: http://localhost:8080/api
- Health Check: http://localhost:8080/health

## API Endpoints

### Drivers
- `POST /api/drivers` - Create a new driver
- `GET /api/drivers` - List all drivers

### Vehicles
- `POST /api/vehicles` - Create a new vehicle
- `GET /api/vehicles` - List all vehicles
- `GET /api/vehicles/{id}/location` - Get vehicle location

### Routes
- `POST /api/routes` - Create a new route
- `GET /api/routes` - List all routes
- `GET /api/routes/{id}` - Get route details
- `PATCH /api/routes/{id}/status` - Update route status
- `POST /api/routes/{id}/optimize` - Optimize route

## Project Structure

```
Go-Dispatch/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database initialization
│   ├── models/            # Data models
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   └── repositories/      # Data access layer
├── web/                   # Frontend files
│   └── index.html        # Single-page application
├── .env.example          # Environment variables template
├── go.mod                # Go module dependencies
└── README.md            # This file
```

## Features Overview

### Driver Management
- Add drivers with contact information
- Track driver status (active, inactive, on_route)
- Assign drivers to routes

### Vehicle Management
- Manage fleet vehicles
- Track vehicle status (available, in_use, maintenance)
- Real-time location tracking
- Vehicle assignment to routes

### Route Planning
- Create routes with multiple stops
- Automatic route optimization using Google Maps
- Calculate total distance and estimated time
- Support for pickup and delivery stops
- Route status tracking (planned, in_progress, completed, cancelled)

### Dashboard
- Real-time statistics
- Fleet status overview
- Recent route activity
- Quick access to all features

## Development

### Building
```bash
go build -o go-dispatch main.go
```

### Running Tests
```bash
go test ./...
```

## Configuration

All configuration is managed through environment variables. See `.env.example` for available options.

Required:
- `GOOGLE_MAPS_API_KEY` - Google Maps API key for route optimization

Optional:
- `MOTIVE_API_KEY` - Motive API key for advanced fleet tracking
- `MOTIVE_API_SECRET` - Motive API secret
- `SERVER_PORT` - Server port (default: 8080)
- `DATABASE_PATH` - SQLite database file path (default: ./dispatch.db)

## Security Notes

- Never commit the `.env` file to version control
- Keep your API keys secure
- Use HTTPS in production
- Implement authentication/authorization for production use

## License

MIT License

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Support

For issues and questions, please open an issue on GitHub.