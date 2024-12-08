package entity

import "time"

type Attachment struct {
	ID        int64
	FileName  string
	Path      *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
