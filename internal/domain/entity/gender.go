package entity

type Gender struct {
	value string
}

var (
	UnknownGender = Gender{value: "unknown"}
	MaleGender    = Gender{value: "male"}
	FemaleGender  = Gender{value: "female"}
)

func GenderFromString(text string) (Gender, error) {
	switch text {
	case MaleGender.value:
		return MaleGender, nil
	case FemaleGender.value:
		return FemaleGender, nil
	default:
		return UnknownGender, UnknownValueError
	}
}

func (g Gender) String() string {
	return g.value
}
