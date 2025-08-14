package config

import (
	"os"
)

// Config holds app configuration populated from environment variables.
type Config struct {
	ServiceName string // e.g., "filmclub-api"
	Env         string // e.g., "local" | "dev" | "production"
	HTTPPort    string // e.g., "8080"
	// Add more as you need later:
	// CassandraHosts []string
	// CassandraKeyspace string
	// TMDBAPIKey string
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Load reads env vars (with sensible defaults) and returns a Config.
func Load() Config {
	return Config{
		ServiceName: getenv("APP_SERVICE_NAME", "filmclub-api"),
		Env:         getenv("APP_ENV", "local"),
		HTTPPort:    getenv("APP_HTTP_PORT", "8080"),
	}
}
