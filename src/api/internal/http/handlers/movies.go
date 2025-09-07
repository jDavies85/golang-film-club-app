// internal/handlers/movies.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jDavies85/golang-film-club-app/api/internal/usecase"
)

type MoviesHandler struct {
	search *usecase.SearchMoviesUC
}

func NewMoviesHandler(search *usecase.SearchMoviesUC) *MoviesHandler {
	return &MoviesHandler{search: search}
}

func (h *MoviesHandler) Search(c *gin.Context) {
	q := c.Query("query")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
		return
	}

	page := atoiDefault(c.Query("page"), 1)
	lang := c.Query("lang")

	var includeAdult *bool
	if v := c.Query("adult"); v != "" {
		b := parseBool(v, false)
		includeAdult = &b
	}

	var year, pry *int
	if ys := c.Query("year"); ys != "" {
		if y, err := strconv.Atoi(ys); err == nil {
			year = &y
		}
	}
	if pys := c.Query("pry"); pys != "" {
		if p, err := strconv.Atoi(pys); err == nil {
			pry = &p
		}
	}

	res, err := h.search.Execute(c.Request.Context(), q, page, lang, includeAdult, year, pry)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Response already matches your frontend schema.
	c.Header("Cache-Control", "public, max-age=60")
	c.JSON(http.StatusOK, gin.H{
		"page":         res.Page,
		"totalPages":   res.TotalPages,
		"totalResults": res.TotalResults,
		"results":      res.Results,
	})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return def
}

func parseBool(s string, def bool) bool {
	switch s {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	}
	return def
}
