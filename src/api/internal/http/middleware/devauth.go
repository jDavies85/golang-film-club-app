package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/jDavies85/golang-film-club-app/api/internal/config"
)

type ctxKey string

const userIDKey ctxKey = "currentUserID"

func WithDevAuth(cfg config.Config) gin.HandlerFunc {
	// Only turn on in non-prod
	if cfg.Env == "production" || !cfg.DevAuthEnabled {
		return func(c *gin.Context) { c.Next() } // no-op
	}

	uid, err := gocql.ParseUUID(cfg.DevUserID)
	if err != nil {
		// Fail fast in dev if misconfigured
		panic("APP_DEV_USER_ID is not a valid UUID: " + err.Error())
	}

	return func(c *gin.Context) {
		// Allow overriding via header for convenience (optional)
		if hdr := c.GetHeader("X-Dev-UserID"); hdr != "" {
			if hUID, e := gocql.ParseUUID(hdr); e == nil {
				uid = hUID
			}
		}
		c.Set(string(userIDKey), uid)
		c.Next()
	}
}

// Helper to fetch current user ID in handlers/usecases
func CurrentUserID(c *gin.Context) (gocql.UUID, bool) {
	v, ok := c.Get(string(userIDKey))
	if !ok {
		return gocql.UUID{}, false
	}
	id, ok := v.(gocql.UUID)
	return id, ok
}

// Optional guard middleware for endpoints that require auth (even in dev)
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := CurrentUserID(c); !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
