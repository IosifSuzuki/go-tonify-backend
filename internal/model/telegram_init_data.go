package model

type TelegramInitData struct {
	QueryID             string
	TelegramUserPayload string
	TelegramUser        TelegramUser
	AuthDate            uint
	Hash                string
}

type TelegramUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsPremium    bool   `json:"is_premium"`
}
