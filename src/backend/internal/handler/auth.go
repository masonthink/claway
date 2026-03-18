package handler

import (
	"html"
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

// XLogin handles GET /api/v1/auth/x
// Redirects user to X OAuth authorization page.
// Optional query param: cli_port (for CLI flow, e.g. "19876")
func (h *AuthHandler) XLogin(c echo.Context) error {
	cliPort := c.QueryParam("cli_port")

	authURL, err := h.svc.GetXAuthURL(cliPort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// XCallback handles GET /api/v1/auth/x/callback
// Receives the authorization code from X and completes the OAuth flow.
func (h *AuthHandler) XCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		errorMsg := c.QueryParam("error")
		if errorMsg == "" {
			errorMsg = "missing authorization code"
		}
		return c.HTML(http.StatusBadRequest, authErrorHTML("Authorization failed: "+errorMsg))
	}

	redirectURL, err := h.svc.HandleXCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, authErrorHTML("Login failed: "+userMessage(err)))
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GitHubLogin handles GET /api/v1/auth/github
func (h *AuthHandler) GitHubLogin(c echo.Context) error {
	cliPort := c.QueryParam("cli_port")

	authURL, err := h.svc.GetGitHubAuthURL(cliPort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles GET /api/v1/auth/github/callback
func (h *AuthHandler) GitHubCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		errorMsg := c.QueryParam("error")
		if errorMsg == "" {
			errorMsg = "missing authorization code"
		}
		return c.HTML(http.StatusBadRequest, authErrorHTML("Authorization failed: "+errorMsg))
	}

	redirectURL, err := h.svc.HandleGitHubCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, authErrorHTML("Login failed: "+userMessage(err)))
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GoogleLogin handles GET /api/v1/auth/google
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	cliPort := c.QueryParam("cli_port")

	authURL, err := h.svc.GetGoogleAuthURL(cliPort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles GET /api/v1/auth/google/callback
func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		errorMsg := c.QueryParam("error")
		if errorMsg == "" {
			errorMsg = "missing authorization code"
		}
		return c.HTML(http.StatusBadRequest, authErrorHTML("Authorization failed: "+errorMsg))
	}

	redirectURL, err := h.svc.HandleGoogleCallback(c.Request().Context(), code, state)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, authErrorHTML("Login failed: "+userMessage(err)))
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// OpenClawCallback handles GET /api/v1/auth/openclaw/callback (legacy)
func (h *AuthHandler) OpenClawCallback(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "code query parameter is required"})
	}

	resp, err := h.svc.HandleOpenClawCallback(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetMe handles GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	user, err := h.svc.GetMe(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateAuthSession handles POST /api/v1/auth/session
// Creates a new auth session for agent-based login flows.
func (h *AuthHandler) CreateAuthSession(c echo.Context) error {
	session, authURL, err := h.svc.CreateAuthSession(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"session_id": session.ID,
		"auth_url":   authURL,
		"expires_at": session.ExpiresAt.Format("2006-01-02T15:04:05Z"),
	})
}

// GetAuthSession handles GET /api/v1/auth/session/:sid
// Returns the current status of an auth session.
func (h *AuthHandler) GetAuthSession(c echo.Context) error {
	sid := c.Param("sid")
	if sid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "session id is required"})
	}

	session, err := h.svc.GetAuthSession(c.Request().Context(), sid)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "session not found or expired"})
	}

	resp := map[string]string{
		"status": session.Status,
	}
	if session.Status == "completed" && session.Token != "" {
		resp["token"] = session.Token
	}

	return c.JSON(http.StatusOK, resp)
}

func authErrorHTML(msg string) string {
	safe := html.EscapeString(msg)
	return `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Claway - Auth Error</title>
<style>body{font-family:system-ui;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;background:#0a0a0a;color:#e5e5e5}
.card{text-align:center;padding:2rem;border-radius:12px;border:1px solid #333;max-width:400px}
a{color:#7c8aff}</style></head>
<body><div class="card"><h2>Authentication Error</h2><p>` + safe + `</p><a href="https://claway.cc">Back to Claway</a></div></body></html>`
}
