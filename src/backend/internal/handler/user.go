package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type UserHandler struct {
	svc *service.Service
}

func NewUserHandler(svc *service.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetUserProfile handles GET /api/v1/public/users/:username
func (h *UserHandler) GetUserProfile(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username is required"})
	}

	profile, err := h.svc.GetUserProfile(c.Request().Context(), username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, profile)
}
