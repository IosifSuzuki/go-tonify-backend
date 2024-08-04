package model

type CreateFreelancer struct {
	TelegramRawInitData string  `json:"telegram_raw_init_data"`
	FirstName           string  `json:"first_name"`
	MiddleName          *string `json:"middle_name"`
	LastName            string  `json:"last_name"`
	Gender              string  `json:"gender"`
	Country             string  `json:"country"`
	City                string  `json:"city"`
	CompanyName         string  `json:"company_name"`
	CompanyDescription  string  `json:"company_description"`
}
