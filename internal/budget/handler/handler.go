package handler

import (
	"expent-backend/internal/budget/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListBudgets(c *gin.Context) {
	shared.SuccessResponse(c, "ListBudgets placeholder", nil)
}

func (h *Handler) CreateBudget(c *gin.Context) {
	shared.SuccessResponse(c, "CreateBudget placeholder", nil)
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	shared.SuccessResponse(c, "DeleteBudget placeholder", nil)
}
