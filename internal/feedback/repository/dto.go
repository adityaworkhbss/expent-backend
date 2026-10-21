package repository

import (
	"time"

	"expent-backend/internal/feedback/model"
)

type FeedbackResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Rating    int       `json:"rating"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToFeedbackResponse(f *model.Feedback) FeedbackResponse {
	return FeedbackResponse{
		ID:        f.ID,
		UserID:    f.UserID,
		Rating:    f.Rating,
		Message:   f.Message,
		CreatedAt: f.CreatedAt,
	}
}
