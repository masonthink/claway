package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type CreditHandler struct {
	svc *service.Service
}

func NewCreditHandler(svc *service.Service) *CreditHandler {
	return &CreditHandler{svc: svc}
}

// GetMyCredits handles GET /api/v1/me/credits
func (h *CreditHandler) GetMyCredits(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.svc.GetMyCredits(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetMyContributions handles GET /api/v1/me/contributions
func (h *CreditHandler) GetMyContributions(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	contribs, err := h.svc.GetMyContributions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"contributions": contribs})
}

// PurchasePRD handles POST /api/v1/prd/:id/purchase
func (h *CreditHandler) PurchasePRD(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	prdID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid prd id"})
	}

	if err := h.svc.PurchasePRD(c.Request().Context(), userID, prdID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "PRD purchased successfully"})
}
