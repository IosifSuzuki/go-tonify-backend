package model

type Account struct {
	ID         *int64  `json:"id"`
	TelegramID *int64  `json:"telegram_id"`
	FirstName  string  `json:"first_name"`
	MiddleName *string `json:"middle_name"`
	LastName   string  `json:"last_name"`
	Nickname   *string `json:"nickname"`
	AboutMe    *string `json:"about_me"`
	Gender     Gender  `json:"gender"`
	Country    *string `json:"country"`
	Location   *string `json:"location"`
	CompanyID  *int64  `json:"company_id"`
}
