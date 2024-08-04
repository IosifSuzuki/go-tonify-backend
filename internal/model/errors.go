package model

import "errors"

var (
	ParametersBadRequestError       = errors.New("parameters bad request error")
	TelegramInitDataDecodeError     = errors.New("telegram init data decode error")
	TelegramInitDataValidationError = errors.New("telegram init data validation error")
	DataBaseOperationError          = errors.New("data base operation error")
	CreationJWTError                = errors.New("creation jwt error")
	ClientAlreadyExistsError        = errors.New("client already exists error")
	FreelancerAlreadyExistsError    = errors.New("freelancer already exists error")
)
