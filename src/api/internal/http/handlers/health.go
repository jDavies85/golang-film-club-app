package handlers

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/jDavies85/golang-film-club-app/api/internal/config"
)

type HealthHandler struct {
	cfg config.Config
}

func NewHealthHandler(cfg config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// Health returns a simple JSON payload; handy for Azure health checks and local smoke tests.
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"service":    h.cfg.ServiceName,
		"env":        h.cfg.Env,
		"goVersion":  runtime.Version(),
		"httpPort":   h.cfg.HTTPPort,
		"uptimeHint": "add process start time later if you like",
	})
}
