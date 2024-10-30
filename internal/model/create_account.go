package model

type CreateAccount struct {
	TelegramRawInitData string  `json:"telegram_raw_init_data" binding:"required"`
	FirstName           string  `json:"first_name" binding:"required"`
	MiddleName          *string `json:"middle_name"`
	LastName            string  `json:"last_name" binding:"required"`
	Nickname            *string `json:"nickname"`
	AboutMe             *string `json:"about_me"`
	Gender              Gender  `json:"gender" enums:"male,female,unknown" binding:"required"`
	Country             string  `json:"country" binding:"required"`
	Location            string  `json:"location" binding:"required"`
	CompanyName         *string `json:"company_name"`
	CompanyDescription  *string `json:"company_description"`
}
