package validator

import (
	"TaskHub/services/user"
	"github.com/go-playground/validator/v10"
)

func VerifyUsernameUnique(fl validator.FieldLevel) bool {

	username := fl.Field().Interface().(string)

	// 如果有错误，则说明用户名未被使用，返回true
	if _, err := user.GetUserByUsername(username); err != nil {
		return true
	}

	return false
}

func VerifyEmailUnique(fl validator.FieldLevel) bool {

	email := fl.Field().Interface().(string)

	// 如果有错误，则说明邮箱未被使用，返回true
	if _, err := user.GetUserByEmail(email); err != nil {
		return true
	}

	return false
}
