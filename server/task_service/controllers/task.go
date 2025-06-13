package controllers

import (
	"TaskHub/task_service/models"
	"TaskHub/task_service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateHandler(c *gin.Context) {
	task := &models.CreateTaskRequest{}
	if err := c.ShouldBindJSON(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	creatorID := c.GetUint("user_id")

	if err := services.CreateTask(task, creatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务创建成功"})
}
