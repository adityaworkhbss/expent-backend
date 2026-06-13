package service

import (
	"context"
	"expent-backend/internal/emi/model"
	"expent-backend/internal/emi/repository"
	"fmt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListEmis(ctx context.Context, userID string) ([]model.Emi, error) {
	return s.repo.ListEmis(userID)
}

func (s *Service) CreateEmi(ctx context.Context, emi model.Emi) (*model.Emi, error) {
	return s.repo.CreateEmi(emi)
}

func (s *Service) UpdateEmi(ctx context.Context, userID, emiID string, emi model.Emi) (*model.Emi, error) {
	e, err := s.repo.GetEmiByID(emiID)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("emi not found")
	}
	if e.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.repo.UpdateEmi(emiID, emi)
}

func (s *Service) DeleteEmi(ctx context.Context, userID, emiID string) error {
	e, err := s.repo.GetEmiByID(emiID)
	if err != nil {
		return err
	}
	if e == nil {
		return fmt.Errorf("emi not found")
	}
	if e.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return s.repo.DeleteEmi(emiID)
}
