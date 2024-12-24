package dto

type GetLikedAccounts struct {
	Offset int64 `form:"offset" example:"5"`
	Limit  int64 `form:"limit" example:"10" binding:"required"`
}
