package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type ComputeHandler struct {
	svc *service.Service
}

func NewComputeHandler(svc *service.Service) *ComputeHandler {
	return &ComputeHandler{svc: svc}
}

// GetMyCompute handles GET /api/v1/me/compute
func (h *ComputeHandler) GetMyCompute(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	resp, err := h.svc.GetMyCompute(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetMyIdeaCompute handles GET /api/v1/me/compute/ideas/:id
func (h *ComputeHandler) GetMyIdeaCompute(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	resp, err := h.svc.GetMyIdeaCompute(c.Request().Context(), userID, ideaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetIdeaCompute handles GET /api/v1/ideas/:id/compute
func (h *ComputeHandler) GetIdeaCompute(c echo.Context) error {
	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	resp, err := h.svc.GetIdeaCompute(c.Request().Context(), ideaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetTaskCompute handles GET /api/v1/tasks/:id/compute
func (h *ComputeHandler) GetTaskCompute(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	resp, err := h.svc.GetTaskCompute(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetPlatformCompute handles GET /api/v1/platform/compute
func (h *ComputeHandler) GetPlatformCompute(c echo.Context) error {
	resp, err := h.svc.GetPlatformCompute(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}
