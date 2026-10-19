package model

import "time"

type DeviceToken struct {
	ID        string
	Token     string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}