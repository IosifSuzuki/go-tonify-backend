package dto

type PostMatchAction struct {
	TargetID int64 `json:"target_id" binding:"required,gt=0" example:"2"`
}
