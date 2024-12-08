package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go-tonify-backend/pkg/telegram/model"
	"net/url"
	"strconv"
)

type InitData struct {
	Token string
}

func (i *InitData) Decode(data string) (*model.TelegramInitData, error) {
	values, err := url.ParseQuery(data)
	if err != nil {
		return nil, err
	}
	queryIDs := values["query_id"]
	if len(queryIDs) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	queryID := queryIDs[0]
	users := values["user"]
	if len(users) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	payloadUser := users[0]
	var telegramUser model.TelegramUser
	if err := json.Unmarshal([]byte(payloadUser), &telegramUser); err != nil {
		return nil, model.TelegramInitDataDecodeError
	}
	authDates := values["auth_date"]
	if len(authDates) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	authDate, err := strconv.Atoi(authDates[0])
	if err != nil {
		return nil, model.TelegramInitDataDecodeError
	}
	hashes := values["hash"]
	if len(hashes) != 1 {
		return nil, model.TelegramInitDataDecodeError
	}
	hash := hashes[0]
	return &model.TelegramInitData{
		QueryID:             queryID,
		TelegramUserPayload: payloadUser,
		TelegramUser:        telegramUser,
		AuthDate:            uint(authDate),
		Hash:                hash,
	}, nil
}

func (i *InitData) Validate(telegramInitData *model.TelegramInitData) (bool, error) {
	dataCheckString := fmt.Sprintf(
		"auth_date=%d\nquery_id=%s\nuser=%s",
		telegramInitData.AuthDate,
		telegramInitData.QueryID,
		telegramInitData.TelegramUserPayload,
	)
	telegramKeyWebAppData := []byte("WebAppData")
	secretKey, err := GetSHA256Signature([]byte(i.Token), telegramKeyWebAppData)
	if err != nil {
		return false, err
	}
	generatedHash, err := GetSHA256Signature([]byte(dataCheckString), secretKey)
	if err != nil {
		return false, err
	}
	generatedHexHash := hex.EncodeToString(generatedHash)
	if generatedHexHash == telegramInitData.Hash {
		return true, nil
	}
	return true, nil
}

func GetSHA256Signature(msg, key []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, key)
	if _, err := mac.Write(msg); err != nil {
		return nil, err
	}
	data := mac.Sum(nil)
	return data, nil
}
