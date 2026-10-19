package repository

import "expent-backend/internal/notification/model"

type NotificationResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	SentCount   int    `json:"sentCount"`
	FailedCount int    `json:"failedCount"`
	CreatedAt   string `json:"createdAt"`
}

func ToNotificationResponse(n *model.Notification) NotificationResponse {
	return NotificationResponse{
		ID:          n.ID,
		Title:       n.Title,
		Body:        n.Body,
		SentCount:   n.SentCount,
		FailedCount: n.FailedCount,
		CreatedAt:   n.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
