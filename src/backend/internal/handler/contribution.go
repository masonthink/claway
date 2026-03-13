package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type ContributionHandler struct {
	svc *service.Service
}

func NewContributionHandler(svc *service.Service) *ContributionHandler {
	return &ContributionHandler{svc: svc}
}

// CreateContribution handles POST /api/v1/ideas/:id/contributions
func (h *ContributionHandler) CreateContribution(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	var req service.CreateContributionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	contrib, err := h.svc.CreateContribution(c.Request().Context(), userID, ideaID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusCreated, contrib)
}

// UpdateContribution handles PUT /api/v1/contributions/:id
func (h *ContributionHandler) UpdateContribution(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	contribID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contribution id"})
	}

	var req service.UpdateContributionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	contrib, err := h.svc.UpdateContribution(c.Request().Context(), userID, contribID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, contrib)
}

// SubmitContribution handles POST /api/v1/contributions/:id/submit
func (h *ContributionHandler) SubmitContribution(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	contribID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contribution id"})
	}

	contrib, err := h.svc.SubmitContribution(c.Request().Context(), userID, contribID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, contrib)
}

// ListContributions handles GET /api/v1/ideas/:id/contributions and GET /api/v1/public/ideas/:id/contributions
func (h *ContributionHandler) ListContributions(c echo.Context) error {
	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	contributions, err := h.svc.ListContributions(c.Request().Context(), ideaID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, contributions)
}

// GetContribution handles GET /api/v1/contributions/:id
func (h *ContributionHandler) GetContribution(c echo.Context) error {
	contribID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contribution id"})
	}

	// Try to get user_id if authenticated, default to 0 for public access
	var userID int64
	if uid, ok := c.Get("user_id").(int64); ok {
		userID = uid
	}

	contrib, err := h.svc.GetContribution(c.Request().Context(), userID, contribID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, contrib)
}

// ListMyContributions handles GET /api/v1/me/contributions
func (h *ContributionHandler) ListMyContributions(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	contributions, total, err := h.svc.ListMyContributions(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"contributions": contributions,
		"total":         total,
	})
}

// GetDraftPreview handles GET /api/v1/draft/:contribution_id
func (h *ContributionHandler) GetDraftPreview(c echo.Context) error {
	userID := c.Get("user_id").(int64)
	contribID, err := strconv.ParseInt(c.Param("contribution_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid contribution id"})
	}

	contrib, err := h.svc.GetDraftPreview(c.Request().Context(), userID, contribID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": userMessage(err)})
	}

	return c.JSON(http.StatusOK, contrib)
}
