package handler

import (
	"context"
	"net/http"

	"expent-backend/internal/device_token/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterDeviceToken(c *gin.Context) {
	var req struct {
		Token  string `json:"token" binding:"required"`
		UserID string `json:"userId" binding:"required"`
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
		shared.ErrorResponse(c, http.StatusForbidden, "Cannot register token for another user")
		return
	}

	ctx := context.Background()
	resp, err := h.svc.RegisterDeviceToken(ctx, req.UserID, req.Token)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	shared.SuccessResponse(c, "Device token registered successfully", resp)
}