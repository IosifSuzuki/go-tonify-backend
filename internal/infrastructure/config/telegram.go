package config

import (
	"go-tonify-backend/internal/domain/entity"
	"os"
	"sync"
)

type Telegram struct {
	BotToken   string
	MiniAppURL string
}

var (
	telegramOnce     sync.Once
	telegramError    error
	telegramInstance *Telegram
)

func GetTelegram() (*Telegram, error) {
	telegramOnce.Do(func() {
		var (
			instance Telegram
			ok       bool
		)
		instance.BotToken, ok = os.LookupEnv("TELEGRAM_BOT_TOKEN")
		if !ok {
			telegramError = entity.NilError
			return
		}
		instance.MiniAppURL, ok = os.LookupEnv("TELEGRAM_BOT_MINI_APP_URL")
		if !ok {
			telegramError = entity.NilError
			return
		}
		telegramInstance = &instance
	})
	return telegramInstance, telegramError
}
