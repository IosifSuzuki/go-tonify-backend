package dto

import "errors"

var (
	NilError                            = errors.New("nil error")
	BadRequestError                     = errors.New("malformed request: required parameters are missing or incorrect")
	InternalServerError                 = errors.New("an unexpected error occurred on the server")
	FailProcessRequestError             = errors.New("the request could not be processed due to an error")
	MissingAuthorizationTokenError      = errors.New("authorization token is missing from the request")
	RoleExpectedError                   = errors.New("role expected error")
	MissingAccountIDError               = errors.New("missing account id")
	ModelNotFoundError                  = errors.New("model not found")
	CastTypeError                       = errors.New("failed to cast the value to the specified type")
	ParseValidateTokenError             = errors.New("failed to parse / validate token")
	DuplicateAccountWithTelegramIDError = errors.New("an account with the specified telegram id already exists")
	CreateTaskLimitError                = errors.New("exceeded the maximum task limit")
)
