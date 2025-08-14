package http

import (
	"github.com/gin-gonic/gin"

	"github.com/jDavies85/golang-film-club-app/api/internal/config"
	"github.com/jDavies85/golang-film-club-app/api/internal/http/handlers"
)

// RegisterRoutes wires all HTTP endpoints.
// Pass cfg so handlers can surface version/name/env on /health, etc.
func RegisterRoutes(r *gin.Engine, cfg config.Config) {
	h := handlers.NewHealthHandler(cfg)

	// basic system endpoints
	r.GET("/health", h.Health)

	// future v1 api group:
	// v1 := r.Group("/v1")
	// {
	//   v1.GET("/movies/search", movieHandler.Search)
	//   v1.POST("/clubs", clubHandler.Create)
	//   ...
	// }
}
