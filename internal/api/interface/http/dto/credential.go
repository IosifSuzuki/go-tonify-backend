package dto

type Credential struct {
	TelegramInitData string `json:"telegram_init_data" validate:"required" example:"query_id=AAGcqlFKAAAAAJyqUUp6-Y62&user=%7B%22id%22%3A1246866076%2C%22first_name%22%3A%22Dante%22%2C%22last_name%22%3A%22%22%2C%22username%22%3A%22S_User%22%2C%22language_code%22%3A%22en%22%7D&auth_date=1651689536&hash=de7f6b26aadbd667a36d76d91969ecf6ffec70ffaa40b3e98d20555e2406bfbb"`
}
