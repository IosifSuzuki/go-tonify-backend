package dto

type GetListTask struct {
	AccountID int64 `form:"account_id" example:"355654520" binding:"required"`
	Offset    int64 `form:"offset" example:"5"`
	Limit     int64 `form:"limit" example:"10" binding:"required"`
}
