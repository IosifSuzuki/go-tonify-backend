package dto

type Pagination struct {
	Offset int64 `json:"offset" example:"0"`
	Limit  int64 `json:"limit" example:"100"`
	Total  int64 `json:"total" example:"32"`
	Data   any   `json:"data"`
}
