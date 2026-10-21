package model

import "time"

type Feedback struct {
	ID        string
	UserID    string
	Rating    int
	Message   string
	CreatedAt time.Time
}
