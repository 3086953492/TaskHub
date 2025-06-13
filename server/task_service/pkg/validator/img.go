package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateImagesURL(fl validator.FieldLevel) bool {
	images := fl.Field().Interface().([]string)
	
	if len(images) == 0 {	// 允许没有图片
		return true
	}
	for _, image := range images {
		if !strings.HasPrefix(image, "http") {
			return false
		}
		if !strings.HasSuffix(image, ".jpg") && !strings.HasSuffix(image, ".png") && !strings.HasSuffix(image, ".jpeg") && !strings.HasSuffix(image, ".webp") {
			return false
		}
	}	
	return true
}