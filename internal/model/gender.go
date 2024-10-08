package model

type Gender string

var (
	Male    Gender = "mela"
	Female  Gender = "female"
	Unknown Gender
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
