package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go-tonify-backend/internal/api/interface/http/validator/impl"
	"go-tonify-backend/internal/container"
	"strings"
)

const (
	enumValidateTag     = "enum_validate"
	nicknameValidateTag = "nickname"
)

type HttpValidator interface {
	Register(*validator.Validate) error
	Translate(validator.ValidationErrors) string
}

type httpValidator struct {
	container container.Container
	enTrans   ut.Translator
}

func NewValidator(container container.Container) HttpValidator {
	return &httpValidator{
		container: container,
	}
}

func (h *httpValidator) Register(validate *validator.Validate) error {
	if err := validate.RegisterValidation(enumValidateTag, impl.ValidateEnum); err != nil {
		return err
	}
	if err := validate.RegisterValidation(nicknameValidateTag, impl.ValidateNickname); err != nil {
		return err
	}
	if err := h.addTranslation(validate); err != nil {
		return err
	}
	return nil
}

func (h *httpValidator) Translate(validationErrors validator.ValidationErrors) string {
	var errList = make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		errList = append(errList, e.Translate(h.enTrans))
	}
	return strings.Join(errList, "\n")
}

func (h *httpValidator) addTranslation(validate *validator.Validate) error {
	english := en.New()
	uni := ut.New(english, english)
	h.enTrans, _ = uni.GetTranslator("en")

	err := validate.RegisterTranslation(enumValidateTag, h.enTrans, func(ut ut.Translator) error {
		return ut.Add(enumValidateTag, "{0} must contains value from enum list", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(enumValidateTag, fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	err = validate.RegisterTranslation(nicknameValidateTag, h.enTrans, func(ut ut.Translator) error {
		return ut.Add(nicknameValidateTag, "{0} must starts with @ and contains numerics, alphas, '-', '_'", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(nicknameValidateTag, fe.Field())
		return t
	})
	if err != nil {
		return err
	}
	return nil
}
