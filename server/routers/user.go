package routers

import (
	"TaskHub/controllers"
	"TaskHub/pkg/auth"
	"github.com/gin-gonic/gin"
)

func LoadUserRouters(router *gin.Engine) {
	userRouters := router.Group("/user")
	{
		userRouters.POST("/login", controllers.Login)
		userRouters.POST("/register", controllers.Register)
		userRouters.PUT("/info", auth.AuthMiddleware(), controllers.UpdateUser)
	}
}