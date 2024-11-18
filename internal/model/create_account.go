package model

type CreateAccount struct {
	TelegramRawInitData string
	FirstName           string
	MiddleName          *string
	LastName            string
	Nickname            *string
	AboutMe             *string
	Gender              Gender
	Country             string
	Location            string
	CompanyID           *int64
	AvatarURL           *string
	DocumentURL         *string
}
