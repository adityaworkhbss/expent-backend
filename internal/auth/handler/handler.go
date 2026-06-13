package handler

import (
	"context"
	"net/http"

	"expent-backend/internal/auth/service"
	"expent-backend/internal/shared"

	"github.com/gin-gonic/gin"
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
	access, refresh, obs, err := h.svc.HandleGoogleLogin(ctx, req.IDToken)
	if err != nil {
		shared.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	shared.SuccessResponse(c, "Logged in", gin.H{"accessToken": access, "refreshToken": refresh, "onboardingCount": obs})
}

// TestLogin handles POST /auth/test-login – bypasses Google verification for local testing.
func (h *Handler) TestLogin(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload or missing email")
		return
	}
	ctx := context.Background()
	access, refresh, obs, err := h.svc.TestLogin(ctx, req.Email)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	shared.SuccessResponse(c, "Test logged in", gin.H{"accessToken": access, "refreshToken": refresh, "onboardingCount": obs})
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

// IncrementOnboarding handles POST /auth/onboarding/increment - updates user onboarding step.
func (h *Handler) IncrementOnboarding(c *gin.Context) {
	var req struct {
		OnboardingCount int `json:"onboardingCount" binding:"required"`
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

	ctx := context.Background()
	err := h.svc.UpdateOnboardingStep(ctx, userId.(string), req.OnboardingCount)
	if err != nil {
		shared.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
