package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/claway/server/internal/service"
)

type DocumentHandler struct {
	svc *service.Service
}

func NewDocumentHandler(svc *service.Service) *DocumentHandler {
	return &DocumentHandler{svc: svc}
}

// GetDocument handles GET /api/v1/tasks/:id/document
func (h *DocumentHandler) GetDocument(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	doc, err := h.svc.GetDocument(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, doc)
}

// ListVersions handles GET /api/v1/tasks/:id/document/versions
func (h *DocumentHandler) ListVersions(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	versions, err := h.svc.ListDocumentVersions(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"versions": versions})
}

// GetVersion handles GET /api/v1/tasks/:id/document/versions/:ver
func (h *DocumentHandler) GetVersion(c echo.Context) error {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	ver, err := strconv.Atoi(c.Param("ver"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid version number"})
	}

	version, err := h.svc.GetDocumentVersion(c.Request().Context(), taskID, ver)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, version)
}

// UpdateDocument handles PUT /api/v1/tasks/:id/document
func (h *DocumentHandler) UpdateDocument(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task id"})
	}

	var req service.UpdateDocumentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := h.svc.UpdateDocument(c.Request().Context(), taskID, userID, req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "document updated"})
}

// PublishPRD handles POST /api/v1/ideas/:id/publish
func (h *DocumentHandler) PublishPRD(c echo.Context) error {
	userID := c.Get("user_id").(int64)

	ideaID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid idea id"})
	}

	prd, err := h.svc.PublishPRD(c.Request().Context(), ideaID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, prd)
}
