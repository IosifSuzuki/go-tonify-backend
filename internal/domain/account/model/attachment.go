package model

import (
	"time"
)

type Attachment struct {
	ID        int64
	Name      string
	Path      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
