package handler

import (
	"expent-backend/internal/dashboard/service"
	"expent-backend/internal/shared"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetDashboard(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		shared.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	dashboard, err := h.svc.GetDashboard(c.Request.Context(), userID.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, "Failed to load dashboard")
		return
	}

	shared.SuccessResponse(c, "Dashboard loaded", dashboard)
}
