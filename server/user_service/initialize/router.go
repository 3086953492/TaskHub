package initialize

import (
	"TaskHub/user_service/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	// 注册路由
	routers.LoadUserRouters(router)
	routers.LoadJwtRouters(router)
	return router
}
