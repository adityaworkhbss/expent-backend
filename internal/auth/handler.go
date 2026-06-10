package auth

import (
    "net/http"
    "context"

    "github.com/gin-gonic/gin"
    "expent-backend/internal/shared"
    "expent-backend/internal/auth/service"
)

// Handler holds dependencies for auth endpoints.
type Handler struct {
    svc *service.Service
}

// NewHandler creates a new auth handler.
func NewHandler(svc *service.Service) *Handler {
    return &Handler{svc: svc}
}

// GoogleLogin handles POST /auth/google – verifies Google ID token and returns JWTs.
func (h *Handler) GoogleLogin(c *gin.Context) {
    var req struct {
        IDToken string `json:"idToken" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
        return
    }
    ctx := context.Background()
    access, refresh, err := h.svc.HandleGoogleLogin(ctx, req.IDToken)
    if err != nil {
        shared.ErrorResponse(c, http.StatusUnauthorized, err.Error())
        return
    }
    shared.SuccessResponse(c, "Logged in", gin.H{"accessToken": access, "refreshToken": refresh})
}

// RefreshToken handles POST /auth/refresh – validates refresh token and issues new access token.
func (h *Handler) RefreshToken(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refreshToken" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
        return
    }
    ctx := context.Background()
    access, err := h.svc.HandleRefresh(ctx, req.RefreshToken)
    if err != nil {
        shared.ErrorResponse(c, http.StatusUnauthorized, err.Error())
        return
    }
    shared.SuccessResponse(c, "Token refreshed", gin.H{"accessToken": access})
}
