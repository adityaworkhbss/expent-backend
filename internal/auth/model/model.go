package model

import "time"

type User struct {
	ID             string
	Email          string
	Name           string
	OnboardingStep int
}

type RefreshToken struct {
	ID        string
	Token     string
	ExpiresAt time.Time
	UserID    string
}
