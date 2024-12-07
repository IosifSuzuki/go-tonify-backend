package dto

type Response[T any] struct {
	Response     *T      `json:"response"`
	ErrorMessage *string `json:"error_message"`
}
