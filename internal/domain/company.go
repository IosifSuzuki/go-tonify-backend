package domain

type Company struct {
	ID          *int64
	Name        *string
	Description *string
}

func NewCompany() *Company {
	return &Company{
		ID:          new(int64),
		Name:        new(string),
		Description: new(string),
	}
}
