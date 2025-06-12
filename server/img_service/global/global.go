package global

import (
	"TaskHub/img_service/configs"
	"go.uber.org/zap"
)

var (
	Cfg      *configs.Config
	Logger   *zap.Logger
)
