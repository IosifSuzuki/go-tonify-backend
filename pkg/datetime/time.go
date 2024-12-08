package datetime

import "time"

type FormatLayoutTime string

const (
	ReadableFormatLayoutTime FormatLayoutTime = "2006-01-02 15:04:05"
)

func GetTimeString(time time.Time, formatLayout FormatLayoutTime) string {
	return time.Format(string(formatLayout))
}
