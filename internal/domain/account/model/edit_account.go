package model

import "mime/multipart"

type EditAccount struct {
	ID                 int64
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
	Tags               *[]string
	CategoryIDs        *[]int64
	AvatarFileHeader   *multipart.FileHeader
	DocumentFileHeader *multipart.FileHeader
}

func (e *EditAccount) HasCompany() bool {
	return e.CompanyName != nil
}
