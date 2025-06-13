package routers

import (
	"TaskHub/task_service/controllers"
	"TaskHub/task_service/middleware"

	"github.com/gin-gonic/gin"
)

func LoadTaskRouter(r *gin.Engine) {
	task := r.Group("/task")
	task.POST("", middleware.AuthMiddleware(), controllers.CreateHandler)
}