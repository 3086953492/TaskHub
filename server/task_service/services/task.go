package services

import (
	"TaskHub/task_service/models"
	"TaskHub/task_service/pkg/logger"
	"TaskHub/task_service/repositories"
	"fmt"
)

func CreateTask(task *models.CreateTaskRequest, creatorID uint) error {
	if err := repositories.CreateTask(task, creatorID); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("创建任务成功, 创建人ID: %d，任务: %v", creatorID, task))
	return nil
}