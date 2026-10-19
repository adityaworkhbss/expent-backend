package handler

import (
	"context"
	"errors"
	"net/http"

	"expent-backend/internal/admin/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ctx := context.Background()
	accessToken, adminResp, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		shared.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	shared.SuccessResponse(c, "Logged in", gin.H{
		"accessToken": accessToken,
		"admin":       adminResp,
	})
}

// Setup creates the first (and only) admin account. It fails once an admin
// already exists, so it cannot be used again after the initial bootstrap.
func (h *Handler) Setup(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: email, password (min 8 chars) are required")
		return
	}

	ctx := context.Background()
	adminResp, err := h.svc.SetupAdmin(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, service.ErrAdminAlreadyExists) {
			shared.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	shared.SuccessResponse(c, "Admin account created", adminResp)
}
