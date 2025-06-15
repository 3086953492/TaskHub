package initialize

import (
	"TaskHub/task_service/global"
	validator_pkg "TaskHub/task_service/pkg/validator"
	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	global.Validate = validator.New()

	if err := global.Validate.RegisterValidation("imageURL", validator_pkg.ValidateImageURL); err != nil {
		return err
	}

	return nil
}
