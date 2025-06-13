package initialize

import (
	"TaskHub/task_service/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	routers.LoadTaskRouters(router)
	return router
}
