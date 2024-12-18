package model

type Message struct {
	ID   int64   `json:"message_id"`
	From *User   `json:"from"`
	Text *string `json:"text"`
	Chat *Chat   `json:"chat"`
}
