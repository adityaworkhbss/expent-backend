package model

import "time"

// User represents a user record from the database.
type User struct {
	ID    string
	Email string
	Name  string
}

// RefreshToken represents a stored refresh token.
type RefreshToken struct {
	ID        string
	Token     string
	ExpiresAt time.Time // using string for simplicity; can be time.Time.
	UserID    string
}

