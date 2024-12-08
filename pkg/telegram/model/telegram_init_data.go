package model

type TelegramInitData struct {
	QueryID             string
	TelegramUserPayload string
	TelegramUser        TelegramUser
	AuthDate            uint
	Hash                string
}
