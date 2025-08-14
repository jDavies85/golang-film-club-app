package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/jDavies85/golang-film-club-app/api/internal/config"
	httprouter "github.com/jDavies85/golang-film-club-app/api/internal/http"
)

func main() {
	// Load .env in local dev; no-op in prod if file is missing.
	_ = godotenv.Load()

	cfg := config.Load()

	// In production, consider gin.SetMode(gin.ReleaseMode)
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Register all routes/handlers
	httprouter.RegisterRoutes(r, cfg)

	addr := ":" + cfg.HTTPPort
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed to start on %s: %v", addr, err)
	}
}
