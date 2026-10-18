package model

import "expent-backend/internal/shared"

// ParseRequest is the input to parse a bank SMS / raw text.
// today, categories, and accounts are optional: the backend loads catalogs
// from the auth token. today should be sent so relative dates use the user's day.
type ParseRequest struct {
	Text       string        `json:"text"`
	RawText    string        `json:"rawText"`
	Today      string        `json:"today"`
	Categories []CatalogItem `json:"categories"`
	Accounts   []CatalogItem `json:"accounts"`
}

// CatalogItem is a user category or account the model may pick from.
type CatalogItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// GeminiParsed is the JSON object the model is asked to return.
type GeminiParsed struct {
	Type          string             `json:"type"`
	Amount        shared.FlexFloat64 `json:"amount"`
	CategoryID    string             `json:"category_id"`
	CategoryName  string             `json:"category_name"`
	Category      string             `json:"category"`
	AccountID     *string            `json:"account_id"`
	AccountName   string             `json:"account_name"`
	PaymentMethod string             `json:"payment_method"`
	Description   string             `json:"description"`
	Merchant      string             `json:"merchant"`
	Note          string             `json:"note"`
	Date          string             `json:"date"`
	Currency      string             `json:"currency"`
	Confidence    shared.FlexFloat64 `json:"confidence"`
}

// ParsedTransaction is returned in data for the app to review then save.
type ParsedTransaction struct {
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	AccountID     *string `json:"account_id"`
	AccountName   string  `json:"account_name"`
	PaymentMethod string  `json:"payment_method"`
	Description   string  `json:"description"`
	Merchant      string  `json:"merchant"`
	Note          string  `json:"note"`
	Date          string  `json:"date"`
	Currency      string  `json:"currency"`
	Confidence    float64 `json:"confidence"`
}

// ParseResult is the full parse-transaction payload.
type ParseResult struct {
	RequiresUserInput     bool              `json:"requires_user_input"`
	MissingRequiredFields []string          `json:"missing_required_fields"`
	Data                  ParsedTransaction `json:"data"`
}
