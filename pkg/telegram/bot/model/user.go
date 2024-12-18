package model

type User struct {
	ID           int64   `json:"id"`
	IsBot        bool    `json:"is_bot"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	Username     *string `json:"username,omitempty"`
	LanguageCode *string `json:"language_code,omitempty"`
}
