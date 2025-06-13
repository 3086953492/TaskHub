package global

import (
	"TaskHub/task_service/configs"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Cfg      *configs.Config
	DB       *gorm.DB
	Validate *validator.Validate
	Logger   *zap.Logger
)
