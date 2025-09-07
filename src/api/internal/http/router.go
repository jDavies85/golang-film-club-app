package http

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jDavies85/golang-film-club-app/api/internal/config"
	"github.com/jDavies85/golang-film-club-app/api/internal/http/handlers"
	"github.com/jDavies85/golang-film-club-app/api/internal/http/middleware"
	"github.com/jDavies85/golang-film-club-app/api/internal/repository/cassandra"
	"github.com/jDavies85/golang-film-club-app/api/internal/usecase"

	// NEW imports for TMDB wiring
	tmdb "github.com/cyruzin/golang-tmdb"
	tmdbadapter "github.com/jDavies85/golang-film-club-app/api/internal/adapters/tmdb"
	moviehandlers "github.com/jDavies85/golang-film-club-app/api/internal/http/handlers"
	movieusecase "github.com/jDavies85/golang-film-club-app/api/internal/usecase"
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

	// ---- TMDB wiring (adapter -> usecase -> handler) ----
	if cfg.TMDBAPIKey == "" {
		log.Println("[WARN] TMDB_API_KEY not set; /v1/movies/search will return 502 errors.")
	}
	tmdbAPI, err := tmdb.Init(cfg.TMDBAPIKey)
	if err != nil {
		// Don't kill the whole server in dev; log and continue.
		log.Printf("[ERROR] tmdb init failed: %v", err)
	}

	// Optional: enforce timeouts on outbound calls
	// If you want, create a custom *http.Client and set via tmdbAPI.SetClient(client)

	const imgBase = "https://image.tmdb.org/t/p"
	adapter := tmdbadapter.New(tmdbAPI, imgBase, "w342", "w780")

	// Default language/enabled-adult fit your app; adjust as needed.
	searchUC := movieusecase.NewSearchMoviesUC(adapter, "en-GB", false)
	mh := moviehandlers.NewMoviesHandler(searchUC)

	// -----------------------------------------------------

	v1 := r.Group("/v1")
	{
		// Clubs
		clubs := handlers.NewClubHandler(clubSvc)
		v1.POST("/clubs", middleware.RequireUser(), clubs.Create)

		// Movies (TMDB search)
		v1.GET("/movies/search", mh.Search)
	}

	// (Optional) graceful shutdown: move session into main() and Close() on exit.
	_ = time.Second // keep imports tidy if you don't use time here
}
