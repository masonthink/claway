package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/clawbeach/server/internal/service"
)

type TaskHandler struct {
	svc *service.Service
}

func NewTaskHandler(svc *service.Service) *TaskHandler {
	return &TaskHandler{svc: svc}
}

// ListTasks handles GET /api/v1/ideas/:id/tasks
func (h *TaskHandler) ListTasks(c echo.Context) error {
	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	tasks, err := h.svc.ListTasks(c.Request().Context(), ideaID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tasks": tasks})
}

// GetTask handles GET /api/v1/tasks/:id
func (h *TaskHandler) GetTask(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	task, err := h.svc.GetTask(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, task)
}

// ClaimTask handles POST /api/v1/tasks/:id/claim
func (h *TaskHandler) ClaimTask(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	if err := h.svc.ClaimTask(c.Request().Context(), taskID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task claimed"})
}

// UnclaimTask handles DELETE /api/v1/tasks/:id/claim
func (h *TaskHandler) UnclaimTask(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	if err := h.svc.UnclaimTask(c.Request().Context(), taskID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task unclaimed"})
}

// SubmitTask handles POST /api/v1/tasks/:id/submit
func (h *TaskHandler) SubmitTask(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	var req service.SubmitTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.svc.SubmitTask(c.Request().Context(), taskID, userID, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task submitted"})
}

// ReviewTask handles POST /api/v1/tasks/:id/review
func (h *TaskHandler) ReviewTask(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	var req service.ReviewTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.svc.ReviewTask(c.Request().Context(), taskID, userID, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task reviewed"})
}
