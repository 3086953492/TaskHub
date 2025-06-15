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
	assigneeIDStr := c.Query("assignee_id")
	creatorIDStr := c.Query("creator_id")

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var page uint
	var pageSize uint
	var assigneeID uint
	var creatorID uint

	if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页码"})
		return
	}

	if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的页大小"})
		return
	}

	conditions := make(map[string]interface{})

	if assigneeIDStr != "" {
		if _, err := fmt.Sscanf(assigneeIDStr, "%d", &assigneeID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的分配人ID"})
			return
		}
		if role != "admin" && assigneeID != userID && assigneeID != 0 {	// 如果分配人ID不为0(未分配)，则只有创建者和分配者可以查看
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看其他人的任务"})
			return
		}
		conditions["assignee_id"] = assigneeID
	}

	if creatorIDStr != "" {
		if _, err := fmt.Sscanf(creatorIDStr, "%d", &creatorID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的创建人ID"})
			return
		}
		if role != "admin" && creatorID != userID {	// 如果创建人ID不为0，则只有创建者可以查看
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看其他人的任务"})
			return
		}
		conditions["creator_id"] = creatorID
	}
	
	if len(conditions) == 0 && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限查看任务列表"})
		return
	}

	tasks, err := services.GetTaskList(page, pageSize, conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func DetailHandler(c *gin.Context) {
	taskIDStr := c.Param("id")

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


func HistoryHandler(c *gin.Context) {
	taskIDStr := c.Param("id")

	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(400, gin.H{"error": "无效的任务ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	history, err := services.GetTaskHistory(taskID, userID, role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, history)
}

func UpdateHandler(c *gin.Context) {

	taskIDStr := c.Param("id")
	var taskID uint
	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	updateTaskRequest := &models.UpdateTaskRequest{}
	if err := c.ShouldBindJSON(updateTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := services.UpdateTask(taskID, userID, role, updateTaskRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务更新成功"})
}

func UpdateStatusHandler(c *gin.Context) {
	taskIDStr := c.Param("id")
	var taskID uint
	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的任务ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	updateTaskStatusRequest := &models.UpdateTaskStatusRequest{}
	if err := c.ShouldBindJSON(updateTaskStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if err := services.UpdateTaskStatus(taskID, userID, role, updateTaskStatusRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "任务状态更新成功"})
}