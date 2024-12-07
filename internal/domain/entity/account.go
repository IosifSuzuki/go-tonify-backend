package entity

import (
	"time"
)

type Account struct {
	ID                   int64
	TelegramID           int64
	FirstName            string
	MiddleName           *string
	LastName             string
	Role                 Role
	Nickname             *string
	AboutMe              *string
	Gender               Gender
	Country              *string
	Location             *string
	CompanyID            *int64
	Company              *Company
	AvatarAttachmentID   *int64
	AvatarAttachment     *Attachment
	DocumentAttachmentID *int64
	DocumentAttachment   *Attachment
	CreatedAt            *time.Time
	UpdatedAt            *time.Time
	DeletedAt            *time.Time
}

func (a *Account) HasCompany() bool {
	return a.CompanyID != nil && a.Company != nil && a.Company.Name != ""
}
