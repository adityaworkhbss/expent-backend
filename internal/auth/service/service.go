package service

import (
	"context"
	"expent-backend/internal/auth/mapper"
	"expent-backend/internal/auth/repository"
	"expent-backend/internal/infrastructure/google"
	"expent-backend/internal/infrastructure/jwt"
	"fmt"
)

type Service struct {
	repo repository.Repository
	// mapper can be used for future DTO conversions
	_ mapper.Mapper
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

// HandleGoogleLogin verifies the Google ID token, ensures a user exists, and returns JWTs.
func (s *Service) HandleGoogleLogin(ctx context.Context, idToken string) (accessToken string, refreshToken string, err error) {
	email, name, err := google.VerifyGoogleIDToken(ctx, idToken)
	if err != nil {
		return "", "", fmt.Errorf("google token validation failed: %w", err)
	}
	// Find or create user
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		user, err = s.repo.CreateUser(email, name)
		if err != nil {
			return "", "", fmt.Errorf("failed to create user: %w", err)
		}
	}
	// Generate JWTs
	accessToken, err = jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err = jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	// Store refresh token in DB
	if err = s.repo.StoreRefreshToken(user.ID, refreshToken); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}
	return accessToken, refreshToken, nil
}

// HandleRefresh validates a refresh token and issues a new access token.
func (s *Service) HandleRefresh(ctx context.Context, refreshToken string) (accessToken string, err error) {
	// Validate refresh token and extract user ID
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}
	// Ensure token exists in DB (optional revocation check)
	stored, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil || stored == nil {
		return "", fmt.Errorf("refresh token not found or revoked")
	}
	// Generate new access token
	accessToken, err = jwt.GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return accessToken, nil
}
