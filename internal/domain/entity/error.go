package entity

import "errors"

var (
	NilError                = errors.New("nil error")
	ConvertStringToIntError = errors.New("convert string to int error")
	EmptyValueError         = errors.New("empty error")
	UnknownValueError       = errors.New("unknown value error")
)
