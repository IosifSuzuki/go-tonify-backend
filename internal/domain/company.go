package domain

import "time"

type Company struct {
	ID          int64
	Name        *string
	Description *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}

func NewCompany() *Company {
	return &Company{
		Name:        new(string),
		Description: new(string),
	}
}
