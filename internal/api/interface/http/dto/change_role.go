package dto

type ChangeRole struct {
	NewRole Role `json:"new_role" binding:"required,enum_validate" example:"client"`
}
