# Go-Dispatch

A fleet dispatch and route management system built with Go, integrating Google Maps API for routing and Motive API for fleet management.

## Features

- **Route Management**: Create, optimize, and track delivery routes
- **Fleet Integration**: Sync vehicles and drivers from Motive
- **Real-time Location**: Track vehicle locations via Motive API
- **Route Optimization**: Optimize waypoint order using Google Maps
- **Distance & Duration**: Calculate accurate route metrics
- **RESTful API**: Complete HTTP API for all operations

## Prerequisites

- Go 1.21 or higher
- Google Maps API key
- Motive API credentials

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

3. Configure environment variables:
```bash
cp .env.example .env
```

Edit `.env` and add your API keys:
```
GOOGLE_MAPS_API_KEY=your_actual_google_maps_api_key
MOTIVE_API_KEY=your_actual_motive_api_key
MOTIVE_API_SECRET=your_actual_motive_api_secret
```

## Usage

### Running the Server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### Building

```bash
go build -o go-dispatch cmd/server/main.go
./go-dispatch
```

## API Endpoints

### Routes

- `POST /api/routes` - Create a new route
- `GET /api/routes` - List all routes
- `GET /api/routes/{id}` - Get a specific route
- `PATCH /api/routes/{id}/status` - Update route status
- `POST /api/routes/optimize` - Optimize route waypoints

### Drivers

- `POST /api/drivers` - Create a new driver
- `GET /api/drivers` - List all drivers
- `POST /api/drivers/sync` - Sync drivers from Motive

### Vehicles

- `POST /api/vehicles` - Create a new vehicle
- `GET /api/vehicles` - List all vehicles
- `GET /api/vehicles/{id}/location` - Get vehicle location
- `POST /api/vehicles/sync` - Sync vehicles from Motive

### Health

- `GET /health` - Health check endpoint

## Example Requests

### Create a Route

```bash
curl -X POST http://localhost:8080/api/routes \
  -H "Content-Type: application/json" \
  -d '{
    "driver_id": "driver-123",
    "vehicle_id": "vehicle-456",
    "stops": [
      {
        "location": {"latitude": 37.7749, "longitude": -122.4194, "address": "San Francisco, CA"},
        "type": "pickup",
        "notes": "Pick up package"
      },
      {
        "location": {"latitude": 37.8044, "longitude": -122.2712, "address": "Oakland, CA"},
        "type": "delivery",
        "notes": "Deliver package"
      }
    ]
  }'
```

### Optimize Route

```bash
curl -X POST http://localhost:8080/api/routes/optimize \
  -H "Content-Type: application/json" \
  -d '{
    "origin": {"latitude": 37.7749, "longitude": -122.4194},
    "destinations": [
      {"latitude": 37.8044, "longitude": -122.2712},
      {"latitude": 37.3382, "longitude": -121.8863},
      {"latitude": 37.4419, "longitude": -122.1430}
    ]
  }'
```

### Sync Vehicles from Motive

```bash
curl -X POST http://localhost:8080/api/vehicles/sync
```

## Project Structure

```
Go-Dispatch/
├── cmd/
│   └── server/           # Main application entry point
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── googlemaps/   # Google Maps API client
│   │   └── motive/       # Motive API client
│   ├── config/           # Configuration management
│   ├── handler/          # HTTP request handlers
│   ├── models/           # Domain models
│   └── service/          # Business logic
├── .env.example          # Environment template
├── .gitignore
├── go.mod
└── README.md
```

## Configuration

All configuration is managed through environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `GOOGLE_MAPS_API_KEY` | Google Maps API key | Yes |
| `MOTIVE_API_KEY` | Motive API key | Yes |
| `MOTIVE_API_SECRET` | Motive API secret | Yes |
| `MOTIVE_BASE_URL` | Motive API base URL | No (defaults to https://api.gomotive.com) |
| `PORT` | Server port | No (defaults to 8080) |
| `NODE_ENV` | Environment (development/production) | No (defaults to development) |

## Development

### Adding New Features

1. Define models in `internal/models/`
2. Implement business logic in `internal/service/`
3. Create HTTP handlers in `internal/handler/`
4. Register routes in `cmd/server/main.go`

### Testing

```bash
go test ./...
```

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.