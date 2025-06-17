package initialize

import (
	"TaskHub/img_service/middleware"
	"TaskHub/img_service/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	// 添加CORS中间件
	router.Use(middleware.CORSMiddleware())

	routers.LoadImgRouters(router)

	return router
}
