package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds the application configuration
type Config struct {
	// Unsplash API settings
	UnsplashAPIKey string
	Timeout        time.Duration
}

// Load reads configuration from a YAML file
func Load(path string) (*Config, error) {
	apiKey := os.Getenv("UNSPLASH_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("UNSPLASH_API_KEY is not set")
	}

	cfg := &Config{
		UnsplashAPIKey: apiKey,
		Timeout:        30 * time.Second,
	}
	return cfg, nil
}
