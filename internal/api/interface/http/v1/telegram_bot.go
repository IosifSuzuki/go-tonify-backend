package v1

import (
	"github.com/gin-gonic/gin"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/pkg/logger"
	"go-tonify-backend/pkg/telegram/bot"
	"go-tonify-backend/pkg/telegram/bot/model"
	"net/http"
)

const (
	telegramBotAvatarPath = "https://tonifyapp-public.s3.eu-central-1.amazonaws.com/thetonifybot-avatar.png"
	markdownParseMode     = "MarkdownV2"
)

type TelegramBotHandler struct {
	container         container.Container
	telegramBotClient bot.Client
}

func NewTelegramBotHandler(container container.Container) *TelegramBotHandler {
	return &TelegramBotHandler{
		container:         container,
		telegramBotClient: bot.NewClient(container.GetTelegramBotToken()),
	}
}

func (t *TelegramBotHandler) Update(ctx *gin.Context) {
	log := t.container.GetLogger()
	update, err := t.telegramBotClient.ParseResponse(ctx.Request.Body)
	if err != nil {
		log.Error("fail to parse request from telegram", logger.FError(err))
		return
	}
	log.Debug("parsed telegram response", logger.F("update", update))
	openAppInlineButton := model.InlineKeyboardButton{
		Text:       "Open app",
		WebAppInfo: &model.WebAppInfo{URL: t.container.GetTelegramMiniAppURL()},
	}
	sendPhoto := model.SendPhoto{
		ChatID:    update.Message.Chat.ID,
		Photo:     telegramBotAvatarPath,
		Caption:   "*Welcome to TONIFY\\!* 🚀\n\nYour gateway to Web3 freelancing\\, powered by decentralised tech\\. Take control of your work and income – secure\\, fair\\, and easy\\!\n\n👇*Tap the button below to start*👇",
		ParseMode: markdownParseMode,
		ReplyMarkup: model.InlineKeyboardMarkup{
			Buttons: [][]model.InlineKeyboardButton{
				{
					openAppInlineButton,
				},
			},
		},
	}
	if err := t.telegramBotClient.Execute(sendPhoto, bot.SendPhotoMethod); err != nil {
		log.Error("fail send response to telegram", logger.FError(err))
		return
	}
	ctx.Status(http.StatusOK)
}
