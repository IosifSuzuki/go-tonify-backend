package dto

type GetMatchAccounts struct {
	Limit int64 `form:"limit" example:"5" binding:"required"`
}
