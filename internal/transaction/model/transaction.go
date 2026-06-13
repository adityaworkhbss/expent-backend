package model

import "time"

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	AccountID   string    `json:"accountId"`
	CategoryID  string    `json:"categoryId"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // e.g., income, expense
	Timestamp   time.Time `json:"timestamp"`
	Date        string    `json:"date"` // ISO 8601 string representation
	Description string    `json:"description,omitempty"`
	Notes       string    `json:"notes,omitempty"` // notes representation
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
