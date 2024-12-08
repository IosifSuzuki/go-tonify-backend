package model

import "mime/multipart"

type CreateAccount struct {
	TelegramInitData   string
	FirstName          string
	MiddleName         *string
	LastName           string
	Role               string
	Nickname           *string
	AboutMe            *string
	Gender             string
	Country            string
	Location           string
	CompanyName        *string
	CompanyDescription *string
	AvatarFileHeader   *multipart.FileHeader
	DocumentFileHeader *multipart.FileHeader
}

func (c *CreateAccount) HasCompany() bool {
	return c.CompanyName != nil && c.CompanyDescription != nil
}
