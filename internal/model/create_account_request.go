package model

import "mime/multipart"

type CreateAccountRequest struct {
	TelegramRawInitData string                `form:"telegram_raw_init_data" binding:"required"`
	FirstName           string                `form:"first_name" binding:"required"`
	MiddleName          *string               `form:"middle_name"`
	LastName            string                `form:"last_name" binding:"required"`
	Nickname            *string               `form:"nickname"`
	AboutMe             *string               `form:"about_me"`
	Gender              Gender                `form:"gender" enums:"male,female,unknown" binding:"required"`
	Country             string                `form:"country" binding:"required"`
	Location            string                `form:"location" binding:"required"`
	CompanyName         *string               `form:"company_name"`
	CompanyDescription  *string               `form:"company_description"`
	Avatar              *multipart.FileHeader `form:"avatar" binding:"required"`
	Document            *multipart.FileHeader `form:"document" binding:"required"`
}
