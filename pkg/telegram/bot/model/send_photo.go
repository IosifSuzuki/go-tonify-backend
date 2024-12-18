package model

type SendPhoto struct {
	ChatID      int64  `json:"chat_id"`
	Photo       string `json:"photo"`
	Caption     string `json:"caption"`
	ParseMode   string `json:"parse_mode"`
	ReplyMarkup any    `json:"reply_markup,omitempty"`
}
