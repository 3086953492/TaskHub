package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateImageURL(fl validator.FieldLevel) bool {
	image := fl.Field().Interface().(string)

	if !strings.HasPrefix(image, "http") {
		return false
	}
	if !strings.HasSuffix(image, ".jpg") && !strings.HasSuffix(image, ".png") && !strings.HasSuffix(image, ".jpeg") && !strings.HasSuffix(image, ".webp") {
		return false
	}
	return true
}
