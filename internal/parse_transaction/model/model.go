package model

// ParseRequest is the input to parse a bank SMS / raw text.
type ParseRequest struct {
	Text    string `json:"text"`
	RawText string `json:"rawText"`
}

// ParsedTransaction is the AI-extracted transaction fields.
type ParsedTransaction struct {
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`        // "expense" | "income"
	Description string  `json:"description"` // merchant / narration
	Date        string  `json:"date"`        // ISO 8601 if detected, else empty
	Currency    string  `json:"currency"`    // e.g. INR, USD
	Confidence  float64 `json:"confidence"`  // 0-1 score from model
}
