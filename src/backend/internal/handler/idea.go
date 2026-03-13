package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusCreated, idea)
}

// ListIdeas handles GET /api/v1/ideas and GET /api/v1/public/ideas
func (h *IdeaHandler) ListIdeas(c echo.Context) error {
	status := c.QueryParam("status")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.svc.ListIdeas(c.Request().Context(), status, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetIdea handles GET /api/v1/ideas/:id and GET /api/v1/public/ideas/:id
func (h *IdeaHandler) GetIdea(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	idea, err := h.svc.GetIdea(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, idea)
}

// ListMyIdeas handles GET /api/v1/me/ideas
func (h *IdeaHandler) ListMyIdeas(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.svc.ListMyIdeas(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetRevealResult handles GET /api/v1/ideas/:id/result and GET /api/v1/public/ideas/:id/result
func (h *IdeaHandler) GetRevealResult(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	result, err := h.svc.GetRevealResult(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, result)
}
