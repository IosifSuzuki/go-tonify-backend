package domain

import (
	"time"
)

type Account struct {
	ID                   int64
	TelegramID           int64
	FirstName            string
	MiddleName           *string
	LastName             string
	Role                 string
	Nickname             *string
	AboutMe              *string
	Gender               string
	Country              *string
	Location             *string
	CompanyID            *int64
	AvatarAttachmentID   *int64
	DocumentAttachmentID *int64
	CreatedAt            *time.Time
	UpdatedAt            *time.Time
	DeletedAt            *time.Time
}

func NewAccount() *Account {
	return &Account{
		MiddleName:           new(string),
		Nickname:             new(string),
		AboutMe:              new(string),
		Country:              new(string),
		Location:             new(string),
		CompanyID:            new(int64),
		AvatarAttachmentID:   new(int64),
		DocumentAttachmentID: new(int64),
		CreatedAt:            new(time.Time),
		UpdatedAt:            new(time.Time),
		DeletedAt:            new(time.Time),
	}
}
