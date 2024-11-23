package model

type CreateAccount struct {
	TelegramRawInitData string
	FirstName           string
	MiddleName          *string
	LastName            string
	Nickname            *string
	Role                Role
	AboutMe             *string
	Gender              Gender
	Country             string
	Location            string
	AvatarID            *int64
	DocumentID          *int64
	CompanyID           *int64
}
