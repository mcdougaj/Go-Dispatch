package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	GoogleMaps GoogleMapsConfig
	Motive     MotiveConfig
	Server     ServerConfig
	Database   DatabaseConfig
}

// GoogleMapsConfig holds Google Maps API configuration
type GoogleMapsConfig struct {
	APIKey string
}

// MotiveConfig holds Motive API configuration
type MotiveConfig struct {
	APIKey    string
	APISecret string
	BaseURL   string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		GoogleMaps: GoogleMapsConfig{
			APIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		},
		Motive: MotiveConfig{
			APIKey:    os.Getenv("MOTIVE_API_KEY"),
			APISecret: os.Getenv("MOTIVE_API_SECRET"),
			BaseURL:   getEnvOrDefault("MOTIVE_BASE_URL", "https://api.gomotive.com"),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "8080"),
			Env:  getEnvOrDefault("NODE_ENV", "development"),
		},
		Database: DatabaseConfig{
			Path: getEnvOrDefault("DATABASE_PATH", "./dispatch.db"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration values are set
func (c *Config) Validate() error {
	if c.GoogleMaps.APIKey == "" || c.GoogleMaps.APIKey == "your_google_maps_api_key_here" {
		return fmt.Errorf("GOOGLE_MAPS_API_KEY is required")
	}
	if c.Motive.APIKey == "" || c.Motive.APIKey == "your_motive_api_key_here" {
		return fmt.Errorf("MOTIVE_API_KEY is required")
	}
	if c.Motive.APISecret == "" || c.Motive.APISecret == "your_motive_api_secret_here" {
		return fmt.Errorf("MOTIVE_API_SECRET is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
