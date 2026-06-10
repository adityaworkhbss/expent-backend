package google

import (
	"context"
	"expent-backend/configs"
	"fmt"

	"google.golang.org/api/idtoken"
)

// VerifyGoogleIDToken validates the ID token issued by Google and returns the user's email and name.
func VerifyGoogleIDToken(ctx context.Context, idToken string) (email string, name string, err error) {
	payload, err := idtoken.Validate(ctx, idToken, configs.AppConfig.GOOGLE_CLIENT_ID)
	if err != nil {
		return "", "", fmt.Errorf("google id token validation failed: %w", err)
	}
	// The payload contains standard claims.
	emailIf, ok := payload.Claims["email"].(string)
	if !ok {
		return "", "", fmt.Errorf("email claim missing in token")
	}
	nameIf, _ := payload.Claims["name"].(string) // name may be optional
	return emailIf, nameIf, nil
}
