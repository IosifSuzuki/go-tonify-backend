package domain

type Client struct {
	ID         *int64
	TelegramID *int64
	FirstName  *string
	MiddleName *string
	LastName   *string
	Gender     *string
	Country    *string
	City       *string
	CompanyID  *int64
}
