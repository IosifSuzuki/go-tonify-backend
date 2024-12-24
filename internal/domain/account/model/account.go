package model

import (
	"go-tonify-backend/internal/domain/category/model"
	"time"
)

type Account struct {
	ID                 int64
	TelegramID         int64
	FirstName          string
	MiddleName         *string
	LastName           string
	Role               string
	Nickname           *string
	AboutMe            *string
	Gender             string
	Country            *string
	Location           *string
	Tags               *[]Tag
	Categories         *[]model.Category
	Company            *Company
	AvatarAttachment   *Attachment
	DocumentAttachment *Attachment
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
