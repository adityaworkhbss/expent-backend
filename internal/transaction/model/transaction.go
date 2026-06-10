package transaction

import "time"

type Transaction struct {
    ID          string    `json:"id"`
    UserID      string    `json:"userId"`
    AccountID   string    `json:"accountId"`
    CategoryID  string    `json:"categoryId"`
    Amount      float64   `json:"amount"`
    Type        string    `json:"type"` // e.g., income, expense
    Timestamp   time.Time `json:"timestamp"`
    Description string    `json:"description,omitempty"`
}
