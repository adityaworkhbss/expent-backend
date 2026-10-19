package service

import (
	"context"

	"expent-backend/internal/device_token/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterDeviceToken(ctx context.Context, userID, token string) (*repository.DeviceTokenResponse, error) {
	dt, err := s.repo.RegisterDeviceToken(userID, token)
	if err != nil {
		return nil, err
	}
	resp := repository.ToDeviceTokenResponse(dt)
	return &resp, nil
}