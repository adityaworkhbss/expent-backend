package model

type UserCustomization struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	Currency string `json:"currency"`
	Theme    string `json:"theme"`
}
