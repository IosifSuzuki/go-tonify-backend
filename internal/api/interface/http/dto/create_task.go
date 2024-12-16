package dto

type CreateTask struct {
	Title       string `json:"title" binding:"required" example:"Tonify"`
	Description string `json:"description" binding:"required" example:"Tonify is a dynamic and innovative company focused on providing cutting-edge solutions to meet the diverse needs of its clients. With a dedicated team of professionals"`
}
