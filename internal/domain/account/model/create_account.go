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
	Tags               *[]string
	CategoryIDs        *[]int64
	CompanyName        *string
	CompanyDescription *string
	AvatarFileHeader   *multipart.FileHeader
	DocumentFileHeader *multipart.FileHeader
}

func (c *CreateAccount) HasCompany() bool {
	return c.CompanyName != nil && c.CompanyDescription != nil
}

func (c *CreateAccount) HasTags() bool {
	if c.Tags == nil {
		return false
	}
	return len(*c.Tags) > 0
}

func (c *CreateAccount) HasCategories() bool {
	if c.CategoryIDs == nil {
		return false
	}
	return len(*c.CategoryIDs) > 0
}
