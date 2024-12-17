package dto

type CreateAccount struct {
	TelegramInitData   string    `form:"telegram_init_data" binding:"required"`
	FirstName          string    `form:"first_name" binding:"required"`
	MiddleName         *string   `form:"middle_name"`
	LastName           string    `form:"last_name" binding:"required"`
	Role               Role      `form:"role" binding:"required,enum_validate"`
	Nickname           string    `form:"nickname" binding:"required,nickname"`
	AboutMe            *string   `form:"about_me"`
	Gender             Gender    `form:"gender" binding:"required,enum_validate"`
	Country            string    `form:"country" binding:"required"`
	Location           string    `form:"location" binding:"required"`
	Tags               *[]string `form:"tags"`
	CompanyName        *string   `form:"company_name"`
	CompanyDescription *string   `form:"company_description"`
}
