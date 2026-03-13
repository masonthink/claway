package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type VoteHandler struct {
	svc *service.Service
}

func NewVoteHandler(svc *service.Service) *VoteHandler {
	return &VoteHandler{svc: svc}
}

// CastVote handles POST /api/v1/ideas/:id/votes
func (h *VoteHandler) CastVote(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	var req service.CastVoteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	vote, err := h.svc.CastVote(c.Request().Context(), userID, ideaID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"voted_at": vote.VotedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// ListMyVotes handles GET /api/v1/me/votes
func (h *VoteHandler) ListMyVotes(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	votes, total, err := h.svc.ListMyVotes(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"votes": votes,
		"total": total,
	})
}
