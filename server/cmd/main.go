package main

import (
	"TaskHub/config"
	"TaskHub/pkg/database"
	"TaskHub/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic("初始化配置失败: " + err.Error())
	}

	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		panic("初始化日志失败: " + err.Error())
	}
	defer logger.Log.Sync()

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		panic("初始化数据库失败: " + err.Error())
	}

	// 初始化Gin路由
	router := gin.Default()
	// 这里可以添加路由配置

	// 启动服务
	port := config.Viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
