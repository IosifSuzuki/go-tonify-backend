package dto

type GetMatchAccounts struct {
	Limit int64 `json:"limit" example:"5" binding:"required,gt=0"`
}
