package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/clawbeach/server/internal/service"
)

type ProxyHandler struct {
	svc *service.Service
}

func NewProxyHandler(svc *service.Service) *ProxyHandler {
	return &ProxyHandler{svc: svc}
}

// Chat handles POST /api/v1/proxy/chat
// Forwards the request to the upstream LLM and records token usage.
// Requires X-Task-ID header to associate usage with a task.
func (h *ProxyHandler) Chat(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskIDStr := c.Request().Header.Get("X-Task-ID")
	if taskIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "X-Task-ID header is required"})
	}

	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid X-Task-ID"})
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "failed to read request body"})
	}

	respBody, err := h.svc.ProxyChat(c.Request().Context(), userID, taskID, body)
	if err != nil {
		// If we got a response body from upstream, return it with the error status
		if respBody != nil {
			return c.JSONBlob(http.StatusBadGateway, respBody)
		}
		return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
	}

	return c.JSONBlob(http.StatusOK, respBody)
}
