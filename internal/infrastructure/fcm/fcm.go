package fcm

import (
	"context"
	"errors"
	"sync"

	"expent-backend/configs"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

var (
	client     *messaging.Client
	clientOnce sync.Once
	clientErr  error
)

// ErrNotConfigured is returned when no Firebase service account credentials
// have been provided. Notification records can still be created without FCM
// configured; sending the actual push simply becomes a no-op error.
var ErrNotConfigured = errors.New("fcm: no Firebase service account credentials configured (set FIREBASE_SERVICE_ACCOUNT_JSON or FIREBASE_SERVICE_ACCOUNT_PATH)")

func getClient(ctx context.Context) (*messaging.Client, error) {
	clientOnce.Do(func() {
		var opts []option.ClientOption
		switch {
		case configs.AppConfig.FIREBASE_SERVICE_ACCOUNT_JSON != "":
			opts = append(opts, option.WithCredentialsJSON([]byte(configs.AppConfig.FIREBASE_SERVICE_ACCOUNT_JSON)))
		case configs.AppConfig.FIREBASE_SERVICE_ACCOUNT_PATH != "":
			opts = append(opts, option.WithCredentialsFile(configs.AppConfig.FIREBASE_SERVICE_ACCOUNT_PATH))
		default:
			clientErr = ErrNotConfigured
			return
		}

		app, err := firebase.NewApp(ctx, nil, opts...)
		if err != nil {
			clientErr = err
			return
		}
		client, clientErr = app.Messaging(ctx)
	})
	return client, clientErr
}

// SendResult reports how many of the given tokens a push was delivered to,
// and which tokens are no longer valid and should be deleted.
type SendResult struct {
	SuccessCount  int
	FailureCount  int
	InvalidTokens []string
}

// SendToTokens sends the same notification to every device token in the
// batch. Invalid/unregistered tokens are reported back so callers can prune
// them from the database.
func SendToTokens(ctx context.Context, tokens []string, title, body string) (SendResult, error) {
	result := SendResult{}
	if len(tokens) == 0 {
		return result, nil
	}

	msgClient, err := getClient(ctx)
	if err != nil {
		return result, err
	}

	const batchSize = 500 // FCM SendEachForMulticast limit
	for start := 0; start < len(tokens); start += batchSize {
		end := start + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		batch := tokens[start:end]

		resp, err := msgClient.SendEachForMulticast(ctx, &messaging.MulticastMessage{
			Tokens: batch,
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
		})
		if err != nil {
			zap.S().Error("fcm: batch send failed", zap.Error(err))
			result.FailureCount += len(batch)
			continue
		}

		result.SuccessCount += resp.SuccessCount
		result.FailureCount += resp.FailureCount
		for i, r := range resp.Responses {
			if r.Success {
				continue
			}
			if messaging.IsRegistrationTokenNotRegistered(r.Error) || messaging.IsInvalidArgument(r.Error) {
				result.InvalidTokens = append(result.InvalidTokens, batch[i])
			}
			zap.S().Warn("fcm: token send failed", zap.String("token", batch[i]), zap.Error(r.Error))
		}
	}

	return result, nil
}
