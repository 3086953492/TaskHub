package validator

import (
	"TaskHub/configs"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	configs.Validate = validator.New()

	if err := configs.Validate.RegisterValidation("usernameUnique", VerifyUsernameUnique); err != nil {
		return err
	}

	if err := configs.Validate.RegisterValidation("emailUnique", VerifyEmailUnique); err != nil {
		return err
	}

	return nil
}
