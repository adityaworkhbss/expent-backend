package service

import (
	"context"
	"strings"

	"expent-backend/internal/feedback/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SubmitFeedback(ctx context.Context, userID string, rating int, message string) (*repository.FeedbackResponse, error) {
	f, err := s.repo.CreateFeedback(userID, rating, strings.TrimSpace(message))
	if err != nil {
		return nil, err
	}
	resp := repository.ToFeedbackResponse(f)
	return &resp, nil
}
