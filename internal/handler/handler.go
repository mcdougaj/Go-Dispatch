package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mcdougaj/Go-Dispatch/internal/models"
	"github.com/mcdougaj/Go-Dispatch/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	dispatch *service.DispatchService
}

// NewHandler creates a new HTTP handler
func NewHandler(dispatch *service.DispatchService) *Handler {
	return &Handler{
		dispatch: dispatch,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Routes
	r.HandleFunc("/api/routes", h.CreateRoute).Methods("POST")
	r.HandleFunc("/api/routes", h.ListRoutes).Methods("GET")
	r.HandleFunc("/api/routes/{id}", h.GetRoute).Methods("GET")
	r.HandleFunc("/api/routes/{id}/status", h.UpdateRouteStatus).Methods("PATCH")
	r.HandleFunc("/api/routes/optimize", h.OptimizeRoute).Methods("POST")

	// Drivers
	r.HandleFunc("/api/drivers", h.CreateDriver).Methods("POST")
	r.HandleFunc("/api/drivers", h.ListDrivers).Methods("GET")
	r.HandleFunc("/api/drivers/sync", h.SyncDrivers).Methods("POST")

	// Vehicles
	r.HandleFunc("/api/vehicles", h.CreateVehicle).Methods("POST")
	r.HandleFunc("/api/vehicles", h.ListVehicles).Methods("GET")
	r.HandleFunc("/api/vehicles/{id}/location", h.GetVehicleLocation).Methods("GET")
	r.HandleFunc("/api/vehicles/sync", h.SyncVehicles).Methods("POST")

	// Health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

// CreateRoute handles POST /api/routes
func (h *Handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var req models.DispatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	route, err := h.dispatch.CreateRoute(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, route)
}

// ListRoutes handles GET /api/routes
func (h *Handler) ListRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.dispatch.ListRoutes(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, routes)
}

// GetRoute handles GET /api/routes/{id}
func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeID := vars["id"]

	route, err := h.dispatch.GetRoute(r.Context(), routeID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, route)
}

// UpdateRouteStatus handles PATCH /api/routes/{id}/status
func (h *Handler) UpdateRouteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	routeID := vars["id"]

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.dispatch.UpdateRouteStatus(r.Context(), routeID, req.Status); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Route status updated"})
}

// OptimizeRoute handles POST /api/routes/optimize
func (h *Handler) OptimizeRoute(w http.ResponseWriter, r *http.Request) {
	var req models.RouteOptimizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.dispatch.OptimizeRoute(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// CreateDriver handles POST /api/drivers
func (h *Handler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var driver models.Driver
	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.dispatch.CreateDriver(r.Context(), &driver); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, driver)
}

// ListDrivers handles GET /api/drivers
func (h *Handler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.dispatch.ListDrivers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, drivers)
}

// SyncDrivers handles POST /api/drivers/sync
func (h *Handler) SyncDrivers(w http.ResponseWriter, r *http.Request) {
	if err := h.dispatch.SyncDriversFromMotive(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Drivers synced successfully"})
}

// CreateVehicle handles POST /api/vehicles
func (h *Handler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.dispatch.CreateVehicle(r.Context(), &vehicle); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, vehicle)
}

// ListVehicles handles GET /api/vehicles
func (h *Handler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.dispatch.ListVehicles(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, vehicles)
}

// GetVehicleLocation handles GET /api/vehicles/{id}/location
func (h *Handler) GetVehicleLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID := vars["id"]

	location, err := h.dispatch.GetVehicleLocation(r.Context(), vehicleID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, location)
}

// SyncVehicles handles POST /api/vehicles/sync
func (h *Handler) SyncVehicles(w http.ResponseWriter, r *http.Request) {
	if err := h.dispatch.SyncVehiclesFromMotive(r.Context()); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Vehicles synced successfully"})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, map[string]string{"error": message})
}
