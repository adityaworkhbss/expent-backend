package handler

import (
	"expent-backend/internal/dashboard/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetDashboard(c *gin.Context) {
	shared.SuccessResponse(c, "GetDashboard placeholder", nil)
}
