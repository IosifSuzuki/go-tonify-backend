package model

import "errors"

var (
	CreateTaskLimitError = errors.New("exceeded the maximum task limit")
	NilError             = errors.New("nil error")
)
