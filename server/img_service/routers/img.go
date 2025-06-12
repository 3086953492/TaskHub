package routers

import (
	"TaskHub/img_service/controllers"
	"TaskHub/img_service/middleware"

	"github.com/gin-gonic/gin"
)

func LoadImgRouters(router *gin.Engine) {
	imgRouters := router.Group("/img")
	{
		imgRouters.POST("/", middleware.AuthMiddleware(), controllers.UploadHandler)
		imgRouters.Static("/", "./img")
	}
}