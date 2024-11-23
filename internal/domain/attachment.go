package domain

import "time"

type Attachment struct {
	ID        int64
	FileName  string
	Path      *string
	Status    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func NewAttachment() *Attachment {
	return &Attachment{
		Path:      new(string),
		CreatedAt: new(time.Time),
		UpdatedAt: new(time.Time),
		DeletedAt: new(time.Time),
	}
}
