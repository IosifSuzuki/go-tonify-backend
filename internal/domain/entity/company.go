package entity

import "time"

type Company struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
