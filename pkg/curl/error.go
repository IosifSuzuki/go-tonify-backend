package curl

import "errors"

var (
	UnsetURLError = errors.New("url is not set")
	ReadBodyError = errors.New("read body error")
)
