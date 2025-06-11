package routers

import (
	"TaskHub/controllers"
	"TaskHub/middleware"
	"github.com/gin-gonic/gin"
)

func LoadUserRouters(router *gin.Engine) {
	userRouters := router.Group("/user")
	{
		userRouters.POST("/login", controllers.LoginHandler)
		userRouters.POST("/register", controllers.RegisterHandler)
		userRouters.PUT("/info", middleware.AuthMiddleware(), controllers.UpdateHandler)
	}
}
