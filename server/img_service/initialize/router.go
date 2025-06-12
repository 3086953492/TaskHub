package initialize

import (
	"TaskHub/img_service/routers"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	routers.LoadImgRouters(router)

	return router
}
