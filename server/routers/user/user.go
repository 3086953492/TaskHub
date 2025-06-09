package user

import (
	"TaskHub/controllers/user"
	"TaskHub/pkg/auth"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/login", user.Login)
		userRouter.POST("/register", user.Register)
		userRouter.POST("/update", auth.AuthMiddleware(), user.UpdateUser)
		userRouter.POST("/logout", auth.AuthMiddleware(), user.Logout)
	}
}