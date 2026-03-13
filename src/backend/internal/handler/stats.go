package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type StatsHandler struct {
	svc *service.Service
}

func NewStatsHandler(svc *service.Service) *StatsHandler {
	return &StatsHandler{svc: svc}
}

// GetPlatformStats handles GET /api/v1/public/stats
func (h *StatsHandler) GetPlatformStats(c echo.Context) error {
	stats, err := h.svc.GetPlatformStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stats)
}
