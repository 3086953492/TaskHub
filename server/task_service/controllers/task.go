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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	creatorID := c.GetUint("user_id")

	if err := services.CreateTask(task, creatorID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "任务创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "任务创建成功",
	})
}

func AssignHandler(c *gin.Context) {

	taskIDStr := c.Query("id")

	// 将字符串转换为uint
	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	userID := c.GetUint("user_id")

	if err := services.AssignTask(taskID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "任务分配失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "任务分配成功",
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的页码",
		})
		return
	}

	if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的页大小",
		})
		return
	}

	conditions := make(map[string]interface{})

	if assigneeIDStr != "" {
		if _, err := fmt.Sscanf(assigneeIDStr, "%d", &assigneeID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "无效的分配人ID",
			})
			return
		}
		if role != "admin" && assigneeID != userID && assigneeID != 0 {	// 如果分配人ID不为0(未分配)，则只有创建者和分配者可以查看
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "无权限查看其他人的任务",
			})
			return
		}
		conditions["assignee_id"] = assigneeID
	}

	if creatorIDStr != "" {
		if _, err := fmt.Sscanf(creatorIDStr, "%d", &creatorID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "无效的创建人ID",
			})
			return
		}
		if role != "admin" && creatorID != userID {	// 如果创建人ID不为0，则只有创建者可以查看
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "无权限查看其他人的任务",
			})
			return
		}
		conditions["creator_id"] = creatorID
	}
	
	if len(conditions) == 0 && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "无权限查看任务列表",
		})
		return
	}

	tasks, err := services.GetTaskList(page, pageSize, conditions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取任务列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取任务列表成功",
		"data": tasks,
	})
}

func DetailHandler(c *gin.Context) {
	taskIDStr := c.Param("id")

	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	task, err := services.GetTaskDetail(taskID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取任务详情失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取任务详情成功",
		"data": task,
	})
}


func HistoryHandler(c *gin.Context) {
	taskIDStr := c.Param("id")

	var taskID uint

	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")
	history, err := services.GetTaskHistory(taskID, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取任务历史失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取任务历史成功",
		"data": history,
	})
}

func UpdateHandler(c *gin.Context) {

	taskIDStr := c.Param("id")
	var taskID uint
	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	updateTaskRequest := &models.UpdateTaskRequest{}
	if err := c.ShouldBindJSON(updateTaskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if err := services.UpdateTask(taskID, userID, role, updateTaskRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "任务更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "任务更新成功",
	})
}

func UpdateStatusHandler(c *gin.Context) {
	taskIDStr := c.Param("id")
	var taskID uint
	if _, err := fmt.Sscanf(taskIDStr, "%d", &taskID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	userID := c.GetUint("user_id")
	role := c.GetString("role")

	updateTaskStatusRequest := &models.UpdateTaskStatusRequest{}
	if err := c.ShouldBindJSON(updateTaskStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if err := services.UpdateTaskStatus(taskID, userID, role, updateTaskStatusRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "任务状态更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "任务状态更新成功",
	})
}