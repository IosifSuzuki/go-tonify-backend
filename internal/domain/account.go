package domain

import "time"

type Account struct {
	ID         *int64
	TelegramID *int64
	FirstName  *string
	MiddleName *string
	LastName   *string
	Nickname   *string
	AboutMe    *string
	Gender     *string
	Country    *string
	Location   *string
	CompanyID  *int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}
