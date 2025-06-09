package routers

import (
	"TaskHub/routers/user"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	user.InitUserRouter(router)

	return router
}
