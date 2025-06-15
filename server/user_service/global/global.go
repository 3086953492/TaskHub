package global

import (
	"TaskHub/user_service/configs"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Cfg      *configs.Config
	DB       *gorm.DB
	Validate *validator.Validate
	Logger   *zap.Logger
	Redis    *redis.Client
)
