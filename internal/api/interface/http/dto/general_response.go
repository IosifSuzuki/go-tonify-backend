package dto

type Response struct {
	Response     any     `json:"response"`
	ErrorMessage *string `json:"error_message" example:"failed to parse / validate token"`
}
