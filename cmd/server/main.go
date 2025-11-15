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
	"github.com/mcdougaj/Go-Dispatch/internal/api/googlemaps"
	"github.com/mcdougaj/Go-Dispatch/internal/api/motive"
	"github.com/mcdougaj/Go-Dispatch/internal/config"
	"github.com/mcdougaj/Go-Dispatch/internal/handler"
	"github.com/mcdougaj/Go-Dispatch/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Google Maps client
	googleMapsClient, err := googlemaps.NewClient(cfg.GoogleMaps.APIKey)
	if err != nil {
		log.Fatalf("Failed to create Google Maps client: %v", err)
	}

	// Initialize Motive client
	motiveClient := motive.NewClient(
		cfg.Motive.BaseURL,
		cfg.Motive.APIKey,
		cfg.Motive.APISecret,
	)

	// Initialize dispatch service
	dispatchService := service.NewDispatchService(googleMapsClient, motiveClient)

	// Initialize HTTP handler
	httpHandler := handler.NewHandler(dispatchService)

	// Setup router
	router := mux.NewRouter()
	httpHandler.RegisterRoutes(router)

	// Add middleware
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting Go-Dispatch server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("%s %s - %v", r.Method, r.RequestURI, time.Since(start))
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
