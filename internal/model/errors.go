package model

import "errors"

var (
	ParametersBadRequestError        = errors.New("parameters bad request error")
	TelegramInitDataDecodeError      = errors.New("telegram init data decode error")
	TelegramInitDataValidationError  = errors.New("telegram init data validation error")
	DataBaseOperationError           = errors.New("data base operation error")
	CreationJWTError                 = errors.New("creation jwt error")
	AccountAlreadyExistsError        = errors.New("account already exists error")
	AccountNotExistsError            = errors.New("account isn't exist in system")
	JWTNotValidError                 = errors.New("jwt isn't valid")
	EmptyAuthorizationHeaderKeyError = errors.New("empty authorization header key")
	MissedAuthorizationTokenError    = errors.New("missed authorization token")
)
