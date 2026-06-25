package google

import (
	"context"
	"expent-backend/configs"
	"fmt"

	"google.golang.org/api/idtoken"
)

func VerifyGoogleIDToken(ctx context.Context, idToken string) (email string, name string, err error) {
	payload, err := idtoken.Validate(ctx, idToken, configs.AppConfig.GOOGLE_CLIENT_ID)
	if err != nil {
		return "", "", fmt.Errorf("google id token validation failed: %w", err)
	}

	emailIf, ok := payload.Claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("email claim missing in token")
	}
	nameIf, ok := payload.Claims["name"].(string)
	if !ok {
		return "", "", fmt.Errorf("name claim missing in token")
	}
	return emailIf, nameIf, nil
}
