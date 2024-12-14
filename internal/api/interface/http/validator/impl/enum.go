package impl

import (
	"github.com/go-playground/validator/v10"
	"go-tonify-backend/internal/api/interface/http/validator/contract"
)

func ValidateEnum(fl validator.FieldLevel) bool {
	field := fl.Field()

	if !field.CanInterface() {
		return true
	}
	if validatable, ok := field.Interface().(contract.EnumValidatable); ok {
		return validatable.Valid()
	}
	return false
}
