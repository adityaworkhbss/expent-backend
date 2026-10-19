package service

import (
	"context"

	deviceTokenRepo "expent-backend/internal/device_token/repository"
	"expent-backend/internal/infrastructure/fcm"
	"expent-backend/internal/notification/repository"
)

type Service struct {
	repo            repository.Repository
	deviceTokenRepo deviceTokenRepo.Repository
}

func NewService(repo repository.Repository, deviceTokenRepo deviceTokenRepo.Repository) *Service {
	return &Service{repo: repo, deviceTokenRepo: deviceTokenRepo}
}

// Broadcast sends a push notification to every registered device token and
// records the notification along with the resulting delivery counts. FCM
// send failures (e.g. credentials not configured) do not block the
// notification record from being created; sentCount/failedCount reflect the
// outcome of the attempt.
func (s *Service) Broadcast(ctx context.Context, title, body, createdBy string) (*repository.NotificationResponse, string, error) {
	tokens, err := s.deviceTokenRepo.GetAllDeviceTokens()
	if err != nil {
		return nil, "", err
	}

	tokenStrings := make([]string, len(tokens))
	for i, t := range tokens {
		tokenStrings[i] = t.Token
	}

	sentCount, failedCount := 0, 0
	pushError := ""
	result, sendErr := fcm.SendToTokens(ctx, tokenStrings, title, body)
	if sendErr == nil {
		sentCount, failedCount = result.SuccessCount, result.FailureCount
		if len(result.InvalidTokens) > 0 {
			_ = s.deviceTokenRepo.DeleteDeviceTokensByTokens(result.InvalidTokens)
		}
	} else {
		failedCount = len(tokenStrings)
		pushError = sendErr.Error()
	}

	n, err := s.repo.CreateNotification(title, body, createdBy, sentCount, failedCount)
	if err != nil {
		return nil, "", err
	}

	resp := repository.ToNotificationResponse(n)
	return &resp, pushError, nil
}
