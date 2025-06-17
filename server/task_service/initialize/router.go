package initialize

import (
	"TaskHub/task_service/middleware"
	"TaskHub/task_service/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

	// 注册路由
	routers.LoadTaskRouters(router)
	return router
}
