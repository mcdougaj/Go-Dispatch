package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcdougaj/go-dispatch/internal/ai"
	"github.com/mcdougaj/go-dispatch/internal/database"
	"github.com/mcdougaj/go-dispatch/internal/maps"
)

// Server represents the API server
type Server struct {
	db          *database.DB
	mapsClient  *maps.Client
	claudeClient *ai.ClaudeClient
	router      *mux.Router
}

// NewServer creates a new API server
func NewServer(db *database.DB, mapsClient *maps.Client, claudeClient *ai.ClaudeClient) *Server {
	s := &Server{
		db:           db,
		mapsClient:   mapsClient,
		claudeClient: claudeClient,
		router:       mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	// CORS middleware
	s.router.Use(corsMiddleware)

	// API routes
	api := s.router.PathPrefix("/api").Subrouter()

	// Fleet data endpoints (GET only)
	api.HandleFunc("/vehicles", s.getVehicles).Methods("GET", "OPTIONS")
	api.HandleFunc("/drivers", s.getDrivers).Methods("GET", "OPTIONS")
	api.HandleFunc("/assets", s.getAssets).Methods("GET", "OPTIONS")
	api.HandleFunc("/geofences", s.getGeofences).Methods("GET", "OPTIONS")

	// AI query endpoint
	api.HandleFunc("/query", s.handleQuery).Methods("POST", "OPTIONS")

	// Maps endpoints
	api.HandleFunc("/geocode", s.handleGeocode).Methods("GET", "OPTIONS")
	api.HandleFunc("/reverse-geocode", s.handleReverseGeocode).Methods("GET", "OPTIONS")
	api.HandleFunc("/distance", s.handleDistance).Methods("GET", "OPTIONS")

	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getVehicles returns all vehicles
func (s *Server) getVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := s.db.GetAllVehicles()
	if err != nil {
		log.Printf("Error getting vehicles: %v", err)
		http.Error(w, "Failed to get vehicles", http.StatusInternalServerError)
		return
	}

	respondJSON(w, vehicles)
}

// getDrivers returns all drivers
func (s *Server) getDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := s.db.GetAllDrivers()
	if err != nil {
		log.Printf("Error getting drivers: %v", err)
		http.Error(w, "Failed to get drivers", http.StatusInternalServerError)
		return
	}

	respondJSON(w, drivers)
}

// getAssets returns all assets
func (s *Server) getAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := s.db.GetAllAssets()
	if err != nil {
		log.Printf("Error getting assets: %v", err)
		http.Error(w, "Failed to get assets", http.StatusInternalServerError)
		return
	}

	respondJSON(w, assets)
}

// getGeofences returns all geofences
func (s *Server) getGeofences(w http.ResponseWriter, r *http.Request) {
	geofences, err := s.db.GetAllGeofences()
	if err != nil {
		log.Printf("Error getting geofences: %v", err)
		http.Error(w, "Failed to get geofences", http.StatusInternalServerError)
		return
	}

	respondJSON(w, geofences)
}

// handleQuery processes natural language queries
func (s *Server) handleQuery(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	// Get current fleet context
	context, err := s.buildFleetContext()
	if err != nil {
		log.Printf("Error building fleet context: %v", err)
		http.Error(w, "Failed to build fleet context", http.StatusInternalServerError)
		return
	}

	// Query Claude
	response, err := s.claudeClient.QueryFleet(req.Query, context)
	if err != nil {
		log.Printf("Error querying Claude: %v", err)
		http.Error(w, "Failed to process query", http.StatusInternalServerError)
		return
	}

	responseTimeMS := int(time.Since(startTime).Milliseconds())

	// Log query
	if err := s.db.LogQuery(req.Query, response.Response, "", nil, responseTimeMS); err != nil {
		log.Printf("Error logging query: %v", err)
	}

	respondJSON(w, map[string]interface{}{
		"query":          req.Query,
		"response":       response.Response,
		"model":          response.Model,
		"response_time_ms": responseTimeMS,
		"input_tokens":   response.InputTokens,
		"output_tokens":  response.OutputTokens,
	})
}

// buildFleetContext creates the context for Claude queries
func (s *Server) buildFleetContext() (ai.FleetContext, error) {
	vehicles, err := s.db.GetAllVehicles()
	if err != nil {
		return ai.FleetContext{}, err
	}

	drivers, err := s.db.GetAllDrivers()
	if err != nil {
		return ai.FleetContext{}, err
	}

	assets, err := s.db.GetAllAssets()
	if err != nil {
		return ai.FleetContext{}, err
	}

	// Build summaries
	var vehicleSummaries []ai.VehicleSummary
	for _, v := range vehicles {
		location := "Unknown"
		if v.CurrentLat != nil && v.CurrentLon != nil {
			// Try to reverse geocode
			if result, err := s.mapsClient.ReverseGeocode(*v.CurrentLat, *v.CurrentLon); err == nil {
				location = result.FormattedAddress
			} else {
				location = fmt.Sprintf("%.4f, %.4f", *v.CurrentLat, *v.CurrentLon)
			}
		}

		status := "Unknown"
		if v.Status != nil {
			status = *v.Status
		}

		vehicleSummaries = append(vehicleSummaries, ai.VehicleSummary{
			VehicleNumber: v.VehicleNumber,
			Status:        status,
			Location:      location,
			Latitude:      v.CurrentLat,
			Longitude:     v.CurrentLon,
			SpeedMPH:      v.SpeedMPH,
			FuelGallons:   v.FuelGallons,
			DriverID:      v.AssignedDriverID,
		})
	}

	var driverSummaries []ai.DriverSummary
	for _, d := range drivers {
		location := "Unknown"
		if d.CurrentLat != nil && d.CurrentLon != nil {
			location = fmt.Sprintf("%.4f, %.4f", *d.CurrentLat, *d.CurrentLon)
		}

		status := "Unknown"
		if d.HOSStatus != nil {
			status = *d.HOSStatus
		}

		driverSummaries = append(driverSummaries, ai.DriverSummary{
			DriverID:         d.DriverID,
			Status:           status,
			Location:         location,
			Latitude:         d.CurrentLat,
			Longitude:        d.CurrentLon,
			HOSRemainingMins: d.HOSRemainingMins,
		})
	}

	var assetSummaries []ai.AssetSummary
	for _, a := range assets {
		location := "Unknown"
		if a.CurrentLat != nil && a.CurrentLon != nil {
			location = fmt.Sprintf("%.4f, %.4f", *a.CurrentLat, *a.CurrentLon)
		}

		assetType := "Unknown"
		if a.AssetType != nil {
			assetType = *a.AssetType
		}

		status := "Unknown"
		if a.Status != nil {
			status = *a.Status
		}

		assetSummaries = append(assetSummaries, ai.AssetSummary{
			AssetName: a.AssetName,
			AssetType: assetType,
			Status:    status,
			Location:  location,
			Latitude:  a.CurrentLat,
			Longitude: a.CurrentLon,
		})
	}

	return ai.FleetContext{
		QueryTime:    time.Now().Format(time.RFC3339),
		VehicleCount: len(vehicles),
		DriverCount:  len(drivers),
		AssetCount:   len(assets),
		Vehicles:     vehicleSummaries,
		Drivers:      driverSummaries,
		Assets:       assetSummaries,
	}, nil
}

// handleGeocode geocodes an address
func (s *Server) handleGeocode(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}

	result, err := s.mapsClient.Geocode(address)
	if err != nil {
		log.Printf("Geocoding error: %v", err)
		http.Error(w, "Geocoding failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, result)
}

// handleReverseGeocode reverse geocodes coordinates
func (s *Server) handleReverseGeocode(w http.ResponseWriter, r *http.Request) {
	var lat, lon float64
	if _, err := fmt.Sscanf(r.URL.Query().Get("lat"), "%f", &lat); err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscanf(r.URL.Query().Get("lon"), "%f", &lon); err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}

	result, err := s.mapsClient.ReverseGeocode(lat, lon)
	if err != nil {
		log.Printf("Reverse geocoding error: %v", err)
		http.Error(w, "Reverse geocoding failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, result)
}

// handleDistance calculates distance between two points
func (s *Server) handleDistance(w http.ResponseWriter, r *http.Request) {
	var originLat, originLon, destLat, destLon float64

	if _, err := fmt.Sscanf(r.URL.Query().Get("origin_lat"), "%f", &originLat); err != nil {
		http.Error(w, "Invalid origin latitude", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscanf(r.URL.Query().Get("origin_lon"), "%f", &originLon); err != nil {
		http.Error(w, "Invalid origin longitude", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscanf(r.URL.Query().Get("dest_lat"), "%f", &destLat); err != nil {
		http.Error(w, "Invalid destination latitude", http.StatusBadRequest)
		return
	}
	if _, err := fmt.Sscanf(r.URL.Query().Get("dest_lon"), "%f", &destLon); err != nil {
		http.Error(w, "Invalid destination longitude", http.StatusBadRequest)
		return
	}

	result, err := s.mapsClient.CalculateDistance(originLat, originLon, destLat, destLon)
	if err != nil {
		log.Printf("Distance calculation error: %v", err)
		http.Error(w, "Distance calculation failed", http.StatusInternalServerError)
		return
	}

	respondJSON(w, result)
}

// handleHealth returns health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
