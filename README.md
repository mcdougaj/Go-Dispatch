# Go-Dispatch

**AI-Powered Fleet Tracking & Dispatch Application**

A modern fleet management system that integrates Motive fleet data with Google Maps tracking and Claude AI for natural language queries.

## 🚀 Features

- **Real-time Fleet Tracking**: View vehicle, driver, and asset locations on an interactive map
- **AI-Powered Queries**: Ask questions about your fleet in natural language using Claude AI
- **Automatic Data Sync**: 30-second refresh from Motive API
- **Google Maps Integration**: Geocoding, reverse geocoding, and distance calculations
- **Modern UI**: React-based dashboard with real-time updates

## 📋 Requirements

### API Keys Required

1. **Motive API Key** - For fleet data
2. **Google Maps API Key** - For mapping and geocoding
3. **Anthropic API Key** - For Claude AI natural language queries

### System Requirements

- Go 1.21+
- Node.js 18+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

## 🛠️ Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/mcdougaj/Go-Dispatch.git
cd Go-Dispatch
```

### 2. Backend Setup

#### Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` and add your API keys:

```env
# Motive API Configuration
MOTIVE_API_KEY=your_motive_api_key_here
MOTIVE_BASE_URL=https://api.gomotive.com

# Google Maps API
GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here

# Claude API (Anthropic)
ANTHROPIC_API_KEY=your_anthropic_api_key_here

# Database
DATABASE_URL=postgres://dispatch_user:dispatch_pass@localhost:5432/motive_dispatch?sslmode=disable

# Redis Cache
REDIS_URL=redis://localhost:6379/0

# Server Configuration
PORT=8080
SYNC_INTERVAL_SECONDS=30
```

#### Install Go Dependencies

```bash
go mod download
```

#### Start Database (Docker)

```bash
docker-compose up -d postgres redis
```

Or use Docker Compose for everything:

```bash
docker-compose up
```

#### Run Backend (Development)

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend
```

#### Install Dependencies

```bash
npm install
```

#### Configure Environment

```bash
cp .env.example .env.local
```

Edit `.env.local`:

```env
REACT_APP_GOOGLE_MAPS_API_KEY=your_google_maps_api_key_here
REACT_APP_API_URL=http://localhost:8080
```

#### Run Frontend (Development)

```bash
npm start
```

The frontend will be available at `http://localhost:3000`

## 📚 API Documentation

### Fleet Data Endpoints (GET Only)

- `GET /api/vehicles` - Get all vehicles
- `GET /api/drivers` - Get all drivers
- `GET /api/assets` - Get all assets (trailers/containers)
- `GET /api/geofences` - Get all geofences

### AI Query Endpoint

- `POST /api/query` - Send natural language query

```json
{
  "query": "Where are all my vehicles?"
}
```

### Maps Endpoints

- `GET /api/geocode?address=<address>` - Geocode an address
- `GET /api/reverse-geocode?lat=<lat>&lon=<lon>` - Reverse geocode coordinates
- `GET /api/distance?origin_lat=<lat>&origin_lon=<lon>&dest_lat=<lat>&dest_lon=<lon>` - Calculate distance

### Health Check

- `GET /health` - Health status

## 🗄️ Database Schema

The application uses PostgreSQL with the following main tables:

- `vehicles` - Vehicle locations and status
- `drivers` - Driver information and HOS data
- `assets` - Trailer/container locations
- `facilities` - Company facilities and customer locations
- `geofences` - Geographic boundaries
- `query_log` - AI query history
- `sync_status` - Data synchronization tracking

## 🎯 Usage Examples

### Natural Language Queries

Try asking:

- "Where are all my vehicles?"
- "Which drivers are available?"
- "Show me empty trailers"
- "Find vehicle TO6076"
- "What's the status of driver 8024360?"

### Viewing Fleet on Map

1. Open the application at `http://localhost:3000`
2. The map displays:
   - **Green arrows** - Moving vehicles
   - **Orange arrows** - Stopped vehicles
   - **Yellow arrows** - Idling vehicles
   - **Blue circles** - Drivers
   - **Gray squares** - Assets (trailers)

3. Click any marker for details
4. Use the side panels to view lists and select items

## 🔧 Development

### Project Structure

```
Go-Dispatch/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── ai/                  # Claude AI integration
│   ├── api/                 # REST API handlers
│   ├── database/            # Database operations
│   ├── maps/                # Google Maps integration
│   ├── motive/              # Motive API client
│   └── sync/                # Data sync service
├── frontend/
│   ├── public/
│   └── src/
│       ├── components/      # React components
│       ├── services/        # API services
│       └── App.js           # Main app component
├── sql/
│   └── schema.sql           # Database schema
├── docker-compose.yml       # Docker services
├── Dockerfile               # Go API container
└── README.md
```

### Building for Production

#### Backend

```bash
CGO_ENABLED=0 GOOS=linux go build -o go-dispatch ./cmd/api
```

Or use Docker:

```bash
docker build -t go-dispatch .
```

#### Frontend

```bash
cd frontend
npm run build
```

The build output will be in `frontend/build/`

## 🚢 Deployment

### Docker Deployment

```bash
docker-compose up -d
```

This starts:
- Go API server (port 8080)
- PostgreSQL database (port 5432)
- Redis cache (port 6379)

### Environment Variables for Production

Make sure to set production values for:
- Database connection strings
- API keys (use secrets management)
- CORS settings if frontend is on different domain

## 🔐 Security Notes

- Never commit API keys to git
- Use environment variables or secrets management
- The `.env` file is in `.gitignore`
- Driver names are never exposed (only IDs for privacy)
- Use HTTPS in production
- Implement authentication/authorization as needed

## 📊 Monitoring

The application logs:
- Data sync operations to `sync_status` table
- AI queries to `query_log` table
- API errors to stdout/stderr

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📝 License

[Add your license here]

## 🆘 Support

For issues or questions:
- Create an issue on GitHub
- Check the Motive API documentation
- Review Google Maps API documentation
- Review Claude API documentation

## 🎉 Acknowledgments

- **Motive** - Fleet management API
- **Google Maps** - Mapping and geocoding
- **Anthropic Claude** - AI natural language processing
- **React** - Frontend framework
- **Go** - Backend language