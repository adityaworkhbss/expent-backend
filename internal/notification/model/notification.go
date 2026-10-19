package model

import "time"

type Notification struct {
	ID          string
	Title       string
	Body        string
	CreatedBy   string
	SentCount   int
	FailedCount int
	CreatedAt   time.Time
}
