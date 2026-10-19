package handler

import (
	"context"
	"net/http"

	"expent-backend/internal/notification/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /notifications — admin-only. Broadcasts a push
// notification to every registered device token and records the result.
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: title and body are required")
		return
	}

	adminId, exists := c.Get("adminId")
	if !exists {
		shared.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := context.Background()
	resp, pushError, err := h.svc.Broadcast(ctx, req.Title, req.Body, adminId.(string))
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := "Notification sent"
	if pushError != "" {
		message = "Notification saved, but push delivery failed: " + pushError
	}
	shared.SuccessResponse(c, message, resp)
}
