package model

type Credential struct {
	TelegramRawInitData string `json:"telegram_raw_init_data" validate:"required"`
}
