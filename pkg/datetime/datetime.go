package datetime

import (
	"time"
)

type Datetime time.Time

func (d Datetime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	formatted := t.Format(time.RFC3339)
	jsonStr := "\"" + formatted + "\""
	return []byte(jsonStr), nil
}

func (d *Datetime) UnmarshalJSON(b []byte) error {
	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
		return NotJSONStringError
	}
	b = b[1 : len(b)-1]
	t, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		return FailToParseDatetimeError
	}
	*d = Datetime(t)
	return nil
}
