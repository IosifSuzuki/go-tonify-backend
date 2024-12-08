package dto

type Company struct {
	ID          int64   `json:"id" example:"1"`
	Name        *string `json:"name" example:"Tonify"`
	Description *string `json:"description" example:"Prospective company with 5-10 employees"`
}
