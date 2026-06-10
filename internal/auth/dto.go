package auth

// DTO definitions for Auth module.

type GoogleLoginRequest struct {
    IDToken string `json:"idToken" binding:"required"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refreshToken" binding:"required"`
}

type AuthResponse struct {
    AccessToken  string `json:"accessToken"`
    RefreshToken string `json:"refreshToken,omitempty"`
}
