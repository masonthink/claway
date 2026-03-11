package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/clawbeach/server/internal/service"
)

type IdeaHandler struct {
	svc *service.Service
}

func NewIdeaHandler(svc *service.Service) *IdeaHandler {
	return &IdeaHandler{svc: svc}
}

// CreateIdea handles POST /api/v1/ideas
func (h *IdeaHandler) CreateIdea(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	var req service.CreateIdeaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	idea, err := h.svc.CreateIdea(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, idea)
}

// ListIdeas handles GET /api/v1/ideas
func (h *IdeaHandler) ListIdeas(c echo.Context) error {
	status := c.QueryParam("status")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.svc.ListIdeas(c.Request().Context(), status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetIdea handles GET /api/v1/ideas/:id
func (h *IdeaHandler) GetIdea(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	idea, err := h.svc.GetIdea(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, idea)
}

// GetIdeaContext handles GET /api/v1/ideas/:id/context
func (h *IdeaHandler) GetIdeaContext(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	ctx, err := h.svc.GetIdeaContext(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, ctx)
}
