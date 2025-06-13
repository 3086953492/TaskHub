package services

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
	"TaskHub/task_service/pkg/logger"
	"TaskHub/task_service/repositories"
	"errors"
	"fmt"
)

func CreateTask(task *models.CreateTaskRequest, creatorID uint) error {
	if err := global.Validate.Struct(task); err != nil {
		return err
	}
	if err := repositories.CreateTask(task, creatorID); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("创建任务成功, 创建人ID: %d，任务: %v", creatorID, task))
	return nil
}

func AssignTask(taskID, userID uint) error {

	// 确保任务未分配
	assignee_id, err := repositories.GetTaskAssigneeID(taskID)

	if err != nil {
		return errors.New("获取分配状态失败:" + err.Error())
	}

	if assignee_id != 0 {
		return errors.New("任务已分配")
	}

	if err := repositories.AssignTask(taskID, userID); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("任务分配成功, 任务ID: %d，分配人ID: %d", taskID, userID))

	return nil
}

func GetTaskList(page, pageSize uint) ([]models.TaskListResponse, error) {
	return repositories.GetTaskList(int(page), int(pageSize))
}
