package model

type Category struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Color  string `json:"color,omitempty"`
	Icon   string `json:"icon,omitempty"`
}
