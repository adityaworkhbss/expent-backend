package service

import (
	"context"
	"time"
	"expent-backend/internal/transaction/model"
	"expent-backend/internal/transaction/repository"
	"fmt"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListTransactions(ctx context.Context, userID string, from, to *time.Time, page, limit *int) ([]model.Transaction, int, error) {
	return s.repo.ListTransactions(userID, from, to, page, limit)
}

func (s *Service) CreateTransaction(ctx context.Context, tx model.Transaction) (*model.Transaction, error) {
	return s.repo.CreateTransaction(tx)
}

func (s *Service) DeleteTransaction(ctx context.Context, userID, txID string) error {
	// verify ownership
	t, err := s.repo.GetTransactionByID(txID)
	if err != nil {
		return err
	}
	if t == nil {
		return fmt.Errorf("transaction not found")
	}
	if t.UserID != userID {
		return fmt.Errorf("unauthorized")
	}
	return s.repo.DeleteTransaction(txID)
}
