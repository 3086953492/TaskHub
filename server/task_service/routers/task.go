package routers

import (
	"TaskHub/task_service/controllers"
	"TaskHub/task_service/middleware"

	"github.com/gin-gonic/gin"
)

func LoadTaskRouters(r *gin.Engine) {
	task := r.Group("/task")
	task.POST("", middleware.AuthMiddleware(), controllers.CreateHandler)
	task.PATCH("", middleware.AuthMiddleware(), controllers.AssignHandler)
	task.GET("/list", middleware.AuthMiddleware(),middleware.AdminAuthMiddleware(), controllers.ListHandler)
	task.GET("/detail", middleware.AuthMiddleware(), controllers.DetailHandler)
}