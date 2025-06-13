package repositories

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
)

func CreateTask(task *models.CreateTaskRequest, creatorID uint) error {

	taskModel := &models.Task{
		Status:    1,
		CreatorID: creatorID,
	}

	if err := global.DB.Create(taskModel).Error; err != nil {
		return err
	}

	taskInfoModel := &models.TaskInfo{
		TaskID:      taskModel.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
	}

	if err := global.DB.Create(taskInfoModel).Error; err != nil {
		return err
	}

	return nil
}
