package model

type InlineKeyboardMarkup struct {
	Buttons [][]InlineKeyboardButton `json:"inline_keyboard"`
}
