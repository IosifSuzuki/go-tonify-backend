package test

import (
	"testing"
)

func TestSuccessValidateTelegramInitData(t *testing.T) {
	//testInitData := "query_id=AAGcqlFKAAAAAJyqUUp6-Y62&user=%7B%22id%22%3A1246866076%2C%22first_name%22%3A%22Dante%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22S_User%22%2C%22language_code%22%3A%22en%22%7D&auth_date=1651689536&hash=de7f6b26aadbd667a36d76d91969ecf6ffec70ffaa40b3e98d20555e2406bfbb"
	//testBotToken := "5139539316:AAGVhDje2A3mB9yA_7l8-TV8xikC7KcudNk"
	//telegramInitData, err := service.DecodeTelegramInitData(testInitData)
	//if err != nil {
	//	t.Errorf("fail to decode telegran init data: %v", err)
	//}
	//dataCheckString := fmt.Sprintf(
	//	"auth_date=%d\nquery_id=%s\nuser=%s",
	//	telegramInitData.AuthDate,
	//	telegramInitData.QueryID,
	//	telegramInitData.TelegramUserPayload,
	//)
	//secretKey, err := encrypt.GetSHA256Signature([]byte(testBotToken), []byte("WebAppData"))
	//if err != nil {
	//	t.Errorf("fail to get a secret key: %v", err)
	//}
	//generatedHash, err := encrypt.GetSHA256Signature([]byte(dataCheckString), secretKey)
	//if err != nil {
	//	t.Errorf("fail to generate a hash key: %v", err)
	//}
	//
	//generatedHexHash := hex.EncodeToString(generatedHash)
	//if generatedHexHash != telegramInitData.Hash {
	//	t.Errorf("is not euqual hash and generated hash")
	//}
}
