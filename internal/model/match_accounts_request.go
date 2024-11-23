package model

type MatchAccountRequest struct {
	Limit int64 `json:"limit" binding:"required"`
}
