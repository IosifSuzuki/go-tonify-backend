package dto

import "time"

type Attachment struct {
	ID        *int64     `json:"id"`
	Name      *string    `json:"name"`
	Path      *string    `json:"path"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
