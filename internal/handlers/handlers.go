package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcdougaj/Go-Dispatch/internal/models"
	"github.com/mcdougaj/Go-Dispatch/internal/repositories"
	"github.com/mcdougaj/Go-Dispatch/internal/services"
)

type Handler struct {
	driverRepo  *repositories.DriverRepository
	vehicleRepo *repositories.VehicleRepository
	routeRepo   *repositories.RouteRepository
	routeService *services.RouteService
}

func NewHandler(db *sql.DB, routeService *services.RouteService) *Handler {
	return &Handler{
		driverRepo:  repositories.NewDriverRepository(db),
		vehicleRepo: repositories.NewVehicleRepository(db),
		routeRepo:   repositories.NewRouteRepository(db),
		routeService: routeService,
	}
}

// Health check
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Driver handlers
func (h *Handler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	driver, err := h.driverRepo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, driver)
}

func (h *Handler) GetDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := h.driverRepo.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, drivers)
}

// Vehicle handlers
func (h *Handler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req models.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	vehicle, err := h.vehicleRepo.Create(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, vehicle)
}

func (h *Handler) GetVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.vehicleRepo.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, vehicles)
}

func (h *Handler) GetVehicleLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	vehicle, err := h.vehicleRepo.GetByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Vehicle not found")
		return
	}

	location := models.VehicleLocation{
		VehicleID:   vehicle.ID,
		Latitude:    vehicle.CurrentLat,
		Longitude:   vehicle.CurrentLng,
		LastUpdated: vehicle.LastUpdated,
	}

	respondJSON(w, http.StatusOK, location)
}

// Route handlers
func (h *Handler) CreateRoute(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	route, err := h.routeService.CreateRoute(&req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, route)
}

func (h *Handler) GetRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.routeService.GetAllRoutes()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, routes)
}

func (h *Handler) GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid route ID")
		return
	}

	route, err := h.routeService.GetRoute(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Route not found")
		return
	}

	respondJSON(w, http.StatusOK, route)
}

func (h *Handler) UpdateRouteStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid route ID")
		return
	}

	var req models.UpdateRouteStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.routeService.UpdateRouteStatus(id, req.Status); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Status updated successfully"})
}

func (h *Handler) OptimizeRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid route ID")
		return
	}

	if err := h.routeService.OptimizeRoute(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	route, err := h.routeService.GetRoute(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, route)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
