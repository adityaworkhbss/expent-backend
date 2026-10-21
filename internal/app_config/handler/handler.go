package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"expent-backend/internal/app_config/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Get handles GET /admin/app-config — admin-only. Returns the config the app is served.
func (h *Handler) Get(c *gin.Context) {
	resp, err := h.svc.Get(c.Request.Context())
	if err != nil {
		h.fail(c, err)
		return
	}
	shared.SuccessResponse(c, "App config loaded", resp)
}

// Update handles PUT /admin/app-config — admin-only. Validates and publishes a new config.
func (h *Handler) Update(c *gin.Context) {
	var req struct {
		Config      json.RawMessage `json:"config" binding:"required"`
		BaseVersion string          `json:"baseVersion"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: config is required")
		return
	}

	resp, err := h.svc.Update(c.Request.Context(), req.Config, req.BaseVersion)
	if err != nil {
		h.fail(c, err)
		return
	}
	adminId, _ := c.Get("adminId")
	zap.S().Infow("app config published", "adminId", adminId, "version", resp.Version)
	shared.SuccessResponse(c, "App config published", resp)
}

func (h *Handler) fail(c *gin.Context, err error) {
	var validationErr *service.ValidationError
	switch {
	case errors.As(err, &validationErr):
		shared.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, service.ErrVersionConflict):
		shared.ErrorResponse(c, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrNotConfigured):
		shared.ErrorResponse(c, http.StatusServiceUnavailable, err.Error())
	default:
		zap.S().Errorw("app config request failed", "error", err)
		shared.ErrorResponse(c, http.StatusBadGateway, "Couldn't reach config storage: "+err.Error())
	}
}
