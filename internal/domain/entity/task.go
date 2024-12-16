package entity

import "time"

type Task struct {
	ID          int64
	OwnerID     int64
	Title       string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
}
