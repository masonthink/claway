package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type AuthHandler struct {
	svc *service.Service
}

func NewAuthHandler(svc *service.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// OpenClawCallback handles GET /api/v1/auth/openclaw/callback
// Exchanges the authorization code for a token, fetches the user profile,
// creates or finds the user, and returns a session JWT.
func (h *AuthHandler) OpenClawCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code query parameter is required"})
	}

	resp, err := h.svc.HandleOpenClawCallback(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetMe handles GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	user, err := h.svc.GetMe(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
