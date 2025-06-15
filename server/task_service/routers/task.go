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
	task.GET("", middleware.AuthMiddleware(), middleware.AdminAuthMiddleware(), controllers.ListHandler)
	task.GET("/:id", middleware.AuthMiddleware(), controllers.DetailHandler)
	task.GET("/unassignedlist", middleware.AuthMiddleware(), controllers.ListUnassigned)
	task.GET("/history", middleware.AuthMiddleware(), controllers.HistoryHandler)
}
