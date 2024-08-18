package model

type CreateAccount struct {
	TelegramRawInitData string  `json:"telegram_raw_init_data" validate:"required"`
	FirstName           string  `json:"first_name" validate:"required"`
	MiddleName          *string `json:"middle_name" validate:"optional"`
	LastName            string  `json:"last_name" validate:"required"`
	Nickname            *string `json:"nickname" validate:"optional"`
	AboutMe             *string `json:"about_me" validate:"optional"`
	Gender              Gender  `json:"gender" enums:"male,female,unknown" validate:"required"`
	Country             string  `json:"country" validate:"required"`
	Location            string  `json:"location" validate:"required"`
	CompanyName         *string `json:"company_name" validate:"optional"`
	CompanyDescription  *string `json:"company_description" validate:"optional"`
}
