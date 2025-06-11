package initialize

import (
	"TaskHub/routers"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	// 注册路由
	routers.LoadUserRouters(router)

	return router
}
