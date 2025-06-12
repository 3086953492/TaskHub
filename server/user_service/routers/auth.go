package routers

import (
	"TaskHub/user_service/controllers"

	"github.com/gin-gonic/gin"
)

func LoadJwtRouters(router *gin.Engine) {
	authRouters := router.Group("/auth")
	{
		authRouters.POST("/refresh", controllers.RefreshHandler) // 此处不用调用验证中间件，因为刷新时会校验token是否有效
	}
}
