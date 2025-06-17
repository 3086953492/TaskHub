package initialize

import (
	"TaskHub/user_service/middleware"
	"TaskHub/user_service/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 注册路由
	routers.LoadUserRouters(router)
	routers.LoadAuthRouters(router)
	return router
}
