package validator

import (
	"TaskHub/configs"
	"TaskHub/services/user"

	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	configs.Validate = validator.New()

	if err := configs.Validate.RegisterValidation("usernameUnique", user.VerifyUsernameUnique); err != nil {
		return err
	}

	if err := configs.Validate.RegisterValidation("emailUnique", user.VerifyEmailUnique); err != nil {
		return err
	}

	return nil
}
