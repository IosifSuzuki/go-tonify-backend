package model

type Gender string

var (
	Male    Gender = "male"
	Female  Gender = "female"
	Unknown Gender = "unknown"
)

func (g Gender) String() string {
	return string(g)
}

func NewGender(gender string) Gender {
	switch gender {
	case string(Male):
		return Male
	case string(Female):
		return Female
	default:
		return Unknown
	}
}
