package model

import "mime/multipart"

type EditAccountRequest struct {
	FirstName          string                `json:"first_name" form:"first_name" binding:"required"`
	MiddleName         *string               `json:"middle_name" form:"middle_name"`
	LastName           string                `json:"last_name" form:"last_name" binding:"required"`
	Role               Role                  `json:"role" form:"role" enums:"client,freelance" binding:"required"`
	Nickname           *string               `json:"nickname" form:"nickname"`
	AboutMe            *string               `json:"about_me" form:"about_me"`
	Gender             Gender                `json:"gender" form:"gender" enums:"male,female,unknown" binding:"required"`
	Country            string                `json:"country" form:"country" binding:"required"`
	Location           string                `json:"location" form:"location" binding:"required"`
	CompanyName        *string               `json:"company_name" form:"company_name"`
	CompanyDescription *string               `json:"company_description" form:"company_description"`
	Avatar             *multipart.FileHeader `json:"avatar" form:"avatar" binding:"required"`
	Document           *multipart.FileHeader `json:"document" form:"document" binding:"required"`
}
