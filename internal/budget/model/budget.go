package model

type Budget struct {
	ID         string  `json:"id"`
	UserID     string  `json:"userId"`
	CategoryID string  `json:"categoryId"`
	Amount     float64 `json:"amount"`
	Period     string  `json:"period"` // e.g., monthly, yearly
}
