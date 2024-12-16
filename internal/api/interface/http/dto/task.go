package dto

import "go-tonify-backend/pkg/datetime"

type Task struct {
	OwnerID     int64              `json:"owner_id" example:"3458728372"`
	Title       string             `json:"title" example:"Create background/avatar for yt"`
	Description string             `json:"description" example:"I expected a professional, highly talented individual with a strong imagination, capable of transforming ideas into avatars and backgrounds"`
	CreatedAt   *datetime.Datetime `json:"created_at" example:"2024-12-07T19:51:48.130157Z"`
	UpdatedAt   *datetime.Datetime `json:"updated_at" example:"2024-12-07T19:51:48.130157Z"`
}
