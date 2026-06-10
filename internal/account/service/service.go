package account

import (
    "context"
    "expent-backend/internal/account/model"
    "expent-backend/internal/account/repository"
)

type Service struct {
    repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) ListAccounts(ctx context.Context, userID string) ([]model.Account, error) {
    return s.repo.ListAccounts(userID)
}

func (s *Service) CreateAccount(ctx context.Context, userID, name, typ string, balance float64) (*model.Account, error) {
    acc := model.Account{UserID: userID, Name: name, Type: typ, Balance: balance}
    return s.repo.CreateAccount(acc)
}

func (s *Service) DeleteAccount(ctx context.Context, userID, accountID string) error {
    // ownership check
    a, err := s.repo.GetAccountByID(accountID)
    if err != nil {
        return err
    }
    if a == nil {
        return fmt.Errorf("account not found")
    }
    if a.UserID != userID {
        return fmt.Errorf("unauthorized")
    }
    return s.repo.DeleteAccount(accountID)
}
