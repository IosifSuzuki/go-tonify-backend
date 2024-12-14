package validator

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"go-tonify-backend/internal/api/interface/http/dto"
	"go-tonify-backend/internal/container"
	"go-tonify-backend/internal/infrastructure/config"
	"go-tonify-backend/pkg/logger"
	"testing"
	"time"
)

type fakeContainer struct{}

func NewFakeContainer() container.Container {
	return &fakeContainer{}
}

func (f *fakeContainer) GetLogger() logger.Logger {
	return logger.NewLogger(logger.DEV, logger.LevelDebug)
}

func (f *fakeContainer) GetTelegramBotToken() string {
	return ""
}

func (f *fakeContainer) GetAWSConfig() *config.AWS {
	return nil
}

func (f *fakeContainer) GetDBConnection() *sql.DB {
	return nil
}

func (f *fakeContainer) GetJWTSecretKey() string {
	return ""
}

func (f *fakeContainer) GetServerConfig() *config.Server {
	return nil
}

func (f *fakeContainer) GetAccessJWTExpiresIn() time.Duration {
	return 0
}

func (f *fakeContainer) GetRefreshJWTExpiresIn() time.Duration {
	return 0
}

func TestEnumValidation(t *testing.T) {
	var test struct {
		Role   dto.Role   `json:"role" validate:"required,enum_validate"`
		Gender dto.Gender `json:"gender" validate:"required,enum_validate"`
	}
	test.Gender = dto.MaleGender
	test.Role = "ff"
	box := NewFakeContainer()
	v := validator.New()
	testValidator := NewValidator(box)
	t.Run("pass success enum validate", func(t *testing.T) {
		if err := testValidator.Register(v); err != nil {
			t.Error("fail to register validator", logger.FError(err))
		}
		if err := v.Struct(test); err != nil {
			//validationErrors, _ := err.(validator.ValidationErrors)
			t.Error("fail to validate struct",
				logger.FError(err),
			)
		}
	})
}
