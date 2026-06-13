package service

import (
	"context"
	"expent-backend/internal/user_customization/model"
	"expent-backend/internal/user_customization/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCustomization(ctx context.Context, userID string) (*model.UserCustomization, error) {
	return s.repo.GetCustomization(userID)
}

func (s *Service) UpdateCustomization(ctx context.Context, userID string, c model.UserCustomization) (*model.UserCustomization, error) {
	return s.repo.UpdateCustomization(userID, c)
}
