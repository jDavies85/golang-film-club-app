package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"

	"github.com/jDavies85/golang-film-club-app/api/internal/http/middleware"
	"github.com/jDavies85/golang-film-club-app/api/internal/usecase"
)

type ClubHandler struct {
	svc *usecase.ClubService
}

func NewClubHandler(svc *usecase.ClubService) *ClubHandler {
	return &ClubHandler{svc: svc}
}

type createClubReq struct {
	Name      string `json:"name" binding:"required"`
	OwnerName string `json:"ownerName"` // optional for now
}

type createClubRes struct {
	ClubID    gocql.UUID `json:"clubId"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (h *ClubHandler) Create(c *gin.Context) {
	uid, ok := middleware.CurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createClubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.svc.Create(c.Request.Context(), usecase.CreateClubInput{
		OwnerID:   uid,
		Name:      req.Name,
		OwnerName: req.OwnerName,
	})
	if err != nil {
		switch err {
		case usecase.ErrEmptyName:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createClubRes{ClubID: out.ClubID, CreatedAt: out.CreatedAt})
}
