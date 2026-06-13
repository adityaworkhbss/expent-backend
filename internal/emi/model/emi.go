package model

import "time"

type Emi struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // "expense" | "income"
	CategoryID  *string   `json:"categoryId"`
	AccountID   *string   `json:"accountId"`
	Date        time.Time `json:"-"`
	DateStr     string    `json:"date"` // ISO string
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
