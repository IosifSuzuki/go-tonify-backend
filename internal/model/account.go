package model

type Account struct {
	ID         *int64
	TelegramID *int64
	FirstName  string
	MiddleName *string
	LastName   string
	Nickname   *string
	AboutMe    *string
	Gender     Gender
	Country    *string
	Location   *string
	CompanyID  *int64
}
