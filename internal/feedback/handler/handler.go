package handler

import (
	"context"
	"net/http"

	"expent-backend/internal/feedback/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SubmitFeedback(c *gin.Context) {
	var req struct {
		UserID  string `json:"userId" binding:"required"`
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Message string `json:"message" binding:"max=2000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userId, exists := c.Get("userId")
	if !exists {
		shared.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if userId.(string) != req.UserID {
		shared.ErrorResponse(c, http.StatusForbidden, "Cannot submit feedback for another user")
		return
	}

	ctx := context.Background()
	resp, err := h.svc.SubmitFeedback(ctx, req.UserID, req.Rating, req.Message)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	shared.SuccessResponse(c, "Feedback submitted successfully", resp)
}
