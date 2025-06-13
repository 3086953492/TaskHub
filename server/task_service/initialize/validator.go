package initialize

import (
	"TaskHub/task_service/global"


	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	global.Validate = validator.New()



	return nil
}
