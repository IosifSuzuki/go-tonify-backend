package dto

type Gender string

const (
	MaleGender   Gender = "male"
	FemaleGender Gender = "female"
	OtherGender  Gender = "other"
)

func (g Gender) Valid() bool {
	switch g {
	case MaleGender, FemaleGender, OtherGender:
		return true
	default:
		return false
	}
}
