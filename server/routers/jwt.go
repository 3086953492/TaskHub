package routers

import (
	"TaskHub/controllers"

	"github.com/gin-gonic/gin"
)

func LoadJwtRouters(router *gin.Engine)	{
	authRouters := router.Group("/jwt")	
	{
		authRouters.POST("/token", controllers.RefreshHandler)
	}
}