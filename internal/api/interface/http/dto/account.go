package dto

import (
	"go-tonify-backend/pkg/datetime"
)

type Account struct {
	ID                 int64              `json:"id" example:"1"`
	TelegramID         int64              `json:"telegram_id" example:"5443222678"`
	FirstName          string             `json:"first_name" example:"Pavel"`
	MiddleName         *string            `json:"middle_name" example:"Michailovich"`
	LastName           string             `json:"last_name" example:"Melnyk"`
	Role               string             `json:"role" example:"client"`
	Nickname           *string            `json:"nickname" example:"@melnyk"`
	AboutMe            *string            `json:"about_me" example:"like when everything good done."`
	Gender             string             `json:"gender" example:"male"`
	Country            *string            `json:"country" example:"Ukraine"`
	Location           *string            `json:"location" example:"Kyiv"`
	Company            *Company           `json:"company"`
	AvatarAttachment   *Attachment        `json:"avatar_attachment"`
	DocumentAttachment *Attachment        `json:"document_attachment"`
	CreatedAt          *datetime.Datetime `json:"created_at" example:"2024-12-07T19:51:48.130157Z"`
	UpdatedAt          *datetime.Datetime `json:"updated_at" example:"2024-12-07T19:51:48.130157Z"`
}
