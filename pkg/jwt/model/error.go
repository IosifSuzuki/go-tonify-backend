package model

import "errors"

var (
	JWTNotValidError = errors.New("jwt isn't valid")
	CreationJWTError = errors.New("creation jwt error")
)
