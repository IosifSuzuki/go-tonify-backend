package dto

type PostMatchAction struct {
	TargetID int64 `json:"target_id" binding:"required" example:"2"`
}
