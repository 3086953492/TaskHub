package controllers

import (
	"TaskHub/task_service/models"
	"TaskHub/task_service/services"
	"fmt"
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

func AssignHandler(c *gin.Context) {

	taskIDStr := c.Query("id")

	// 将字符串转换为uint
	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID := c.GetUint("user_id")

	if err := services.AssignTask(taskID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务分配成功"})
}

func ListHandler(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	var page uint
	var pageSize uint

	if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页码"})
		return
	}

	if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页大小"})
		return
	}

	tasks, err := services.GetTaskList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func DetailHandler(c *gin.Context) {
	taskIDStr := c.Query("id")

	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	task, err := services.GetTaskDetail(taskID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func ListUnassigned(c *gin.Context) {
	pageStr := c.Query("page")
	pageSizeStr := c.Query("page_size")

	var page uint
	var pageSize uint

	if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
		c.JSON(400, gin.H{"error": "无效的页码"})
	}

	if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
		c.JSON(400, gin.H{"error": "无效的页大小"})
	}

	tasks, err := services.GetUnassignedTasks(page, pageSize)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, tasks)
}

func HistoryHandler(c *gin.Context) {
	taskIDStr := c.Query("id")

	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(400, gin.H{"error": "无效的任务ID"})
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	history, err := services.GetTaskHistory(taskID, userID, role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, history)
}