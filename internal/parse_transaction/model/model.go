package model

// ParseRequest is the input to parse a bank SMS / raw text.
type ParseRequest struct {
	// RawText is the raw bank SMS or transaction description to parse.
	RawText string `json:"rawText" binding:"required"`
}

// ParsedTransaction is the AI-extracted transaction fields.
type ParsedTransaction struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`        // "expense" | "income"
	Description string  `json:"description"` // merchant / narration
	Date        string  `json:"date"`        // ISO 8601 if detected, else empty
	Currency    string  `json:"currency"`    // e.g. INR, USD
	Confidence  float64 `json:"confidence"`  // 0-1 score from model
}
