package dto

type GetCategories struct {
	Offset int64 `form:"offset" example:"0"`
	Limit  int64 `form:"limit" example:"10" binding:"required"`
}
