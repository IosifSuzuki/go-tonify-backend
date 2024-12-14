package impl

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateNickname(fl validator.FieldLevel) bool {
	var re = regexp.MustCompile(`^@[\w\d_-]+$`)
	text := fl.Field().Interface().(string)
	return re.MatchString(text)
}
