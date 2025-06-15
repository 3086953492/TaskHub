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
	task.GET("", middleware.AuthMiddleware(), controllers.ListHandler)
	task.GET("/:id", middleware.AuthMiddleware(), controllers.DetailHandler)
	task.GET("/history/:id", middleware.AuthMiddleware(), controllers.HistoryHandler)
	task.PATCH("/:id", middleware.AuthMiddleware(), controllers.UpdateHandler)
	task.PATCH("/status/:id", middleware.AuthMiddleware(), controllers.UpdateStatusHandler)
}
