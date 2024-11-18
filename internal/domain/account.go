package domain

import "time"

type Account struct {
	ID           *int64
	TelegramID   *int64
	FirstName    *string
	MiddleName   *string
	LastName     *string
	Nickname     *string
	AboutMe      *string
	Gender       *string
	Country      *string
	Location     *string
	CompanyID    *int64
	AvatarPath   *string
	DocumentPath *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
}

func NewAccount() *Account {
	return &Account{
		ID:           new(int64),
		TelegramID:   new(int64),
		FirstName:    new(string),
		MiddleName:   new(string),
		LastName:     new(string),
		Nickname:     new(string),
		AboutMe:      new(string),
		Gender:       new(string),
		Country:      new(string),
		Location:     new(string),
		CompanyID:    new(int64),
		AvatarPath:   new(string),
		DocumentPath: new(string),
		CreatedAt:    new(time.Time),
		UpdatedAt:    new(time.Time),
		DeletedAt:    new(time.Time),
	}
}
