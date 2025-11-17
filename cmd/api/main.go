package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mcdougaj/go-dispatch/internal/ai"
	"github.com/mcdougaj/go-dispatch/internal/api"
	"github.com/mcdougaj/go-dispatch/internal/database"
	"github.com/mcdougaj/go-dispatch/internal/maps"
	"github.com/mcdougaj/go-dispatch/internal/motive"
	"github.com/mcdougaj/go-dispatch/internal/sync"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get configuration from environment
	config := getConfig()

	// Initialize database connection
	log.Println("Connecting to database...")
	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	// Initialize Motive API client
	log.Println("Initializing Motive API client...")
	motiveClient := motive.NewClient(config.MotiveAPIKey, config.MotiveBaseURL)

	// Initialize Google Maps client
	log.Println("Initializing Google Maps client...")
	mapsClient, err := maps.NewClient(config.GoogleMapsAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Google Maps client: %v", err)
	}

	// Initialize Claude AI client
	log.Println("Initializing Claude AI client...")
	claudeClient := ai.NewClaudeClient(config.AnthropicAPIKey)

	// Start sync service
	log.Println("Starting data sync service...")
	syncService := sync.NewService(motiveClient, db, config.SyncIntervalSeconds)
	syncService.Start()
	defer syncService.Stop()

	// Create API server
	log.Println("Starting API server...")
	server := api.NewServer(db, mapsClient, claudeClient)

	// Start HTTP server
	addr := ":" + config.Port
	log.Printf("Server listening on %s", addr)
	log.Println("API endpoints:")
	log.Println("  GET  /api/vehicles")
	log.Println("  GET  /api/drivers")
	log.Println("  GET  /api/assets")
	log.Println("  GET  /api/geofences")
	log.Println("  POST /api/query")
	log.Println("  GET  /api/geocode")
	log.Println("  GET  /api/reverse-geocode")
	log.Println("  GET  /api/distance")
	log.Println("  GET  /health")

	// Setup graceful shutdown
	go func() {
		if err := http.ListenAndServe(addr, server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

// Config holds application configuration
type Config struct {
	Port                string
	DatabaseURL         string
	MotiveAPIKey        string
	MotiveBaseURL       string
	GoogleMapsAPIKey    string
	AnthropicAPIKey     string
	SyncIntervalSeconds int
}

// getConfig reads configuration from environment variables
func getConfig() Config {
	port := getEnv("PORT", "8080")
	syncInterval := getEnv("SYNC_INTERVAL_SECONDS", "30")
	syncIntervalInt, _ := strconv.Atoi(syncInterval)

	return Config{
		Port:                port,
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		MotiveAPIKey:        getEnv("MOTIVE_API_KEY", ""),
		MotiveBaseURL:       getEnv("MOTIVE_BASE_URL", "https://api.gomotive.com"),
		GoogleMapsAPIKey:    getEnv("GOOGLE_MAPS_API_KEY", ""),
		AnthropicAPIKey:     getEnv("ANTHROPIC_API_KEY", ""),
		SyncIntervalSeconds: syncIntervalInt,
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
