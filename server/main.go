package main

import (
	"TaskHub/configs"
	"TaskHub/pkg/auth"
	"TaskHub/pkg/database"
	"TaskHub/pkg/logger"
	"TaskHub/pkg/validator"
	"TaskHub/routers"
	"fmt"
	"log"
	"os"
	"strconv"
	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := configs.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		logger.Error("初始化数据库失败", zap.Error(err))
		return
	}

	// 初始化会话存储
	auth.InitSessionStore()

	// 初始化验证器
	if err := validator.InitValidator(); err != nil {
		logger.Error("初始化验证器失败", zap.Error(err))
		return
	}

	// 获取端口号，优先使用命令行参数
	port := configs.Cfg.Server.Port
	if len(os.Args) > 1 {
		argPort, err := strconv.Atoi(os.Args[1])
		if err == nil {
			port = argPort
		}
	}

	// 初始化 Gin 服务器
	r := routers.InitRouter()

	logger.Info("TaskHub 服务启动成功")
	logger.Info(fmt.Sprintf("服务正在运行在端口 %d", port))

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error("启动服务失败", zap.Error(err))
	}
}
