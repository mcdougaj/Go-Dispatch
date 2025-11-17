package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mcdougaj/Go-Dispatch/internal/config"
	"github.com/mcdougaj/Go-Dispatch/internal/database"
	"github.com/mcdougaj/Go-Dispatch/internal/handlers"
	"github.com/mcdougaj/Go-Dispatch/internal/repositories"
	"github.com/mcdougaj/Go-Dispatch/internal/services"
	"googlemaps.github.io/maps"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database initialized successfully")

	// Initialize Google Maps client
	mapsClient, err := maps.NewClient(maps.WithAPIKey(cfg.GoogleMapsAPIKey))
	if err != nil {
		log.Fatalf("Failed to create Google Maps client: %v", err)
	}

	// Initialize services
	routeRepo := repositories.NewRouteRepository(db)
	routeService := services.NewRouteService(routeRepo, mapsClient)

	// Initialize handlers
	handler := handlers.NewHandler(db, routeService)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/drivers", handler.CreateDriver).Methods("POST")
	api.HandleFunc("/drivers", handler.GetDrivers).Methods("GET")
	api.HandleFunc("/vehicles", handler.CreateVehicle).Methods("POST")
	api.HandleFunc("/vehicles", handler.GetVehicles).Methods("GET")
	api.HandleFunc("/vehicles/{id}/location", handler.GetVehicleLocation).Methods("GET")
	api.HandleFunc("/routes", handler.CreateRoute).Methods("POST")
	api.HandleFunc("/routes", handler.GetRoutes).Methods("GET")
	api.HandleFunc("/routes/{id}", handler.GetRoute).Methods("GET")
	api.HandleFunc("/routes/{id}/status", handler.UpdateRouteStatus).Methods("PATCH")
	api.HandleFunc("/routes/{id}/optimize", handler.OptimizeRoute).Methods("POST")

	// Health check
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Serve static files (frontend)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	// Setup server
	addr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on http://localhost%s", addr)
		log.Printf("Web interface: http://localhost%s", addr)
		log.Printf("API endpoint: http://localhost%s/api", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
