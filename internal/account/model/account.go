package account

type Account struct {
    ID      string  `json:"id"`
    UserID  string  `json:"userId"`
    Name    string  `json:"name"`
    Type    string  `json:"type"`
    Balance float64 `json:"balance,omitempty"`
}
