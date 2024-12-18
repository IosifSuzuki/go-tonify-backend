package model

type InlineKeyboardButton struct {
	Text       string      `json:"text"`
	WebAppInfo *WebAppInfo `json:"web_app,omitempty"`
}
