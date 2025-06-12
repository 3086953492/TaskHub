package initialize

import (
	"TaskHub/user_service/global"
	validator_pkg "TaskHub/user_service/pkg/validator"

	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	global.Validate = validator.New()

	if err := global.Validate.RegisterValidation("usernameUnique", validator_pkg.VerifyUsernameUnique); err != nil {
		return err
	}

	if err := global.Validate.RegisterValidation("emailUnique", validator_pkg.VerifyEmailUnique); err != nil {
		return err
	}

	return nil
}
