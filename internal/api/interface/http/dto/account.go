package dto

import "time"

type Account struct {
	ID                 int64       `json:"id"`
	TelegramID         int64       `json:"telegram_id"`
	FirstName          string      `json:"first_name"`
	MiddleName         *string     `json:"middle_name"`
	LastName           string      `json:"last_name"`
	Role               string      `json:"role"`
	Nickname           *string     `json:"nickname"`
	AboutMe            *string     `json:"about_me"`
	Gender             string      `json:"gender"`
	Country            *string     `json:"country"`
	Location           *string     `json:"location"`
	Company            *Company    `json:"company"`
	AvatarAttachment   *Attachment `json:"avatar_attachment"`
	DocumentAttachment *Attachment `json:"document_attachment"`
	CreatedAt          *time.Time  `json:"created_at"`
	UpdatedAt          *time.Time  `json:"updated_at"`
}
