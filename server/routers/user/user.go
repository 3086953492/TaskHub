package user

import (
	"TaskHub/controllers"
	"TaskHub/pkg/auth"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/login", controllers.Login)
		userRouter.POST("/register", controllers.Register)
		userRouter.POST("/update", auth.AuthMiddleware(), controllers.UpdateUser)
	}
}