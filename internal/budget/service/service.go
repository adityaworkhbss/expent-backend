package service

import (
	"context"
	"expent-backend/internal/budget/model"
	"expent-backend/internal/budget/repository"
	"fmt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListBudgets(ctx context.Context, userID string) ([]model.Budget, error) {
	return s.repo.ListBudgets(userID)
}

func (s *Service) CreateBudget(ctx context.Context, budget model.Budget) (*model.Budget, error) {
	return s.repo.CreateBudget(budget)
}

func (s *Service) UpdateBudget(ctx context.Context, userID, budgetID string, budget model.Budget) (*model.Budget, error) {
	b, err := s.repo.GetBudgetByID(budgetID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, fmt.Errorf("budget not found")
	}
	if b.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.repo.UpdateBudget(budgetID, budget)
}

func (s *Service) DeleteBudget(ctx context.Context, userID, budgetID string) error {
	b, err := s.repo.GetBudgetByID(budgetID)
	if err != nil {
		return err
	}
	if b == nil {
		return fmt.Errorf("budget not found")
	}
	if b.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return s.repo.DeleteBudget(budgetID)
}
