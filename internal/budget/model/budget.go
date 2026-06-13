package model

import "time"

type Budget struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	CategoryID   string    `json:"categoryId"`
	Period       string    `json:"periodType"`  // maps to periodType
	Amount       float64   `json:"limitAmount"` // maps to limitAmount
	StartDate    time.Time `json:"-"`
	EndDate      time.Time `json:"-"`
	StartDateStr string    `json:"startDate"`
	EndDateStr   string    `json:"endDate,omitempty"`
}
