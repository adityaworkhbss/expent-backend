package jwt

import (
	"expent-backend/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed JWT access token for a given user ID.
func GenerateAccessToken(userID string) (string, error) {
	expiresIn, err := time.ParseDuration(configs.AppConfig.JWT_EXPIRES_IN)
	if err != nil {
		expiresIn = 15 * time.Minute // fallback default
	}
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(configs.AppConfig.JWT_SECRET))
	if err != nil {
		zap.S().Error("failed to sign JWT", zap.Error(err))
		return "", err
	}
	return signed, nil
}

// GenerateRefreshToken creates a signed JWT refresh token.
func GenerateRefreshToken(userID string) (string, error) {
	expiresIn, err := time.ParseDuration(configs.AppConfig.JWT_REFRESH_EXPIRES_IN)
	if err != nil {
		expiresIn = 7 * 24 * time.Hour // fallback default
	}
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(configs.AppConfig.JWT_REFRESH_SECRET))
	if err != nil {
		zap.S().Error("failed to sign refresh JWT", zap.Error(err))
		return "", err
	}
	return signed, nil
}

// ValidateAccessToken parses and validates an access token string.
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfig.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// ValidateRefreshToken parses and validates a refresh token string.
func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(configs.AppConfig.JWT_REFRESH_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
