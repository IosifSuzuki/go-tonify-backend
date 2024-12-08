package dto

type Attachment struct {
	ID   *int64  `json:"id" example:"1"`
	Name *string `json:"name" example:"c6141704-d2df-4ada-96b0-01f227109681.png"`
	Path *string `json:"path" example:"https://tonifyapp-attachment.s3.eu-central-1.amazonaws.com/46e31a1d-e2f5-4787-8450-dca5604c4350.png"`
}
