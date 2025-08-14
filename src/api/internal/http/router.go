package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jDavies85/golang-film-club-app/api/internal/config"
	"github.com/jDavies85/golang-film-club-app/api/internal/http/handlers"
	"github.com/jDavies85/golang-film-club-app/api/internal/http/middleware"
	"github.com/jDavies85/golang-film-club-app/api/internal/repository/cassandra"
	"github.com/jDavies85/golang-film-club-app/api/internal/usecase"
)

func RegisterRoutes(r *gin.Engine, cfg config.Config) {
	h := handlers.NewHealthHandler(cfg)
	r.GET("/health", h.Health)

	// Dev auth
	r.Use(middleware.WithDevAuth(cfg))

	// Cassandra session (in real life, create once in main and pass in)
	sess, err := cassandra.NewSession([]string{"127.0.0.1"}, "filmclub", "LOCAL_QUORUM", 5)
	if err != nil {
		panic(err)
	}

	// Repos
	clubRepo := cassandra.NewClubRepo(sess)
	membersRepo := cassandra.NewClubMembersRepo(sess)
	userClubsRepo := cassandra.NewUserClubsRepo(sess)
	guardRepo := cassandra.NewMembershipGuardRepo(sess)

	// Service
	clubSvc := usecase.NewClubService(clubRepo, membersRepo, userClubsRepo, guardRepo)

	v1 := r.Group("/v1")
	{
		clubs := handlers.NewClubHandler(clubSvc)
		v1.POST("/clubs", middleware.RequireUser(), clubs.Create)
	}

	// (Optional) graceful shutdown: move session into main() and Close() on exit.
	_ = time.Second // keep imports tidy if you don't use time here
}
