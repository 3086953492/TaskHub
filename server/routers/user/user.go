package user

import (
	"TaskHub/controllers/user"
	"github.com/gin-gonic/gin"
)

func InitUserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/login", user.Login)
	}
}
