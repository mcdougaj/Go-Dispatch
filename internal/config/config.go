package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	ServerPort       string
	GoogleMapsAPIKey string
	MotiveAPIKey     string
	MotiveAPISecret  string
	DatabasePath     string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (optional)
	_ = godotenv.Load()

	config := &Config{
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		GoogleMapsAPIKey: os.Getenv("GOOGLE_MAPS_API_KEY"),
		MotiveAPIKey:     os.Getenv("MOTIVE_API_KEY"),
		MotiveAPISecret:  os.Getenv("MOTIVE_API_SECRET"),
		DatabasePath:     getEnv("DATABASE_PATH", "./dispatch.db"),
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.GoogleMapsAPIKey == "" {
		return errors.New("GOOGLE_MAPS_API_KEY is required")
	}
	// Motive API credentials are optional
	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
