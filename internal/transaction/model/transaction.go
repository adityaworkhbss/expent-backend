package model

import "time"

// Transaction is one actual income or expense.
// Date is the RFC 3339 timestamp of when it happened; requests may also send a plain yyyy-MM-dd.
type Transaction struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	AccountID  string    `json:"accountId"`
	CategoryID string    `json:"categoryId"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"` // INCOME or EXPENSE
	Date       string    `json:"date"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	// Parsed form of Date, used when writing to the database
	Timestamp time.Time `json:"-"`
}
