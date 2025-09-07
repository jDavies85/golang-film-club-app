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
	TMDBAPIKey     string
	DevAuthEnabled bool   // enable fake auth
	DevUserID      string // UUID string of your seeded user
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
		ServiceName:    getenv("APP_SERVICE_NAME", "filmclub-api"),
		Env:            getenv("APP_ENV", "local"),
		HTTPPort:       getenv("APP_HTTP_PORT", "8080"),
		TMDBAPIKey:     getenv("TMDB_API_KEY", ""),
		DevAuthEnabled: getenv("APP_DEV_AUTH_ENABLED", "true") == "true",
		DevUserID:      getenv("APP_DEV_USER_ID", "12345678-1234-1234-1234-123456789abc"),
	}
}
