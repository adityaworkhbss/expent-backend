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
	DateStr     string    `json:"date"`
	StartDate   string    `json:"startDate"`
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tenure      *int      `json:"tenure,omitempty"`
	MonthsPaid  *int      `json:"monthsPaid,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
