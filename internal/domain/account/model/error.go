package model

import "errors"

var (
	NilError                            = errors.New("nil error")
	DecodeTelegramInitDataError         = errors.New("decode telegram initialization data error")
	InvalidTelegramInitDataError        = errors.New("invalid telegram initialization data provided")
	DuplicateAccountWithTelegramIDError = errors.New("an account with the specified telegram id already exists")
	EntityNotFoundError                 = errors.New("entity not found")
	DuplicateMatchActionError           = errors.New("duplicate match action")
	UnhandledMatchActionError           = errors.New("unhandled match action")
)
