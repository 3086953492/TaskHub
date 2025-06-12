package main

import (
	"TaskHub/img_service/global"
	"TaskHub/img_service/initialize"
	"TaskHub/img_service/pkg/logger"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
)

func main() {
	// 初始化配置
	if err := initialize.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化日志
	if err := initialize.InitLogger(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 获取端口号，优先使用命令行参数
	port := global.Cfg.Server.Port
	if len(os.Args) > 1 {
		argPort, err := strconv.Atoi(os.Args[1])
		if err == nil {
			port = argPort
		}
	}

	// 初始化 Gin 路由
	r := initialize.InitRouters()

	logger.Info("TaskHub/img_service 服务启动成功")
	logger.Info(fmt.Sprintf("服务正在运行在端口 %d", port))

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Error("启动服务失败", zap.Error(err))
	}
}
