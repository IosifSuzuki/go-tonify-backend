package datetime

import "errors"

var (
	NotJSONStringError       = errors.New("not a json string")
	FailToParseDatetimeError = errors.New("failed to parse time")
)
