package repositories

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
	"strconv"
	"time"
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

	for i, image := range task.Images {
		taskImageModel := &models.TaskImage{
			TaskID:    taskModel.ID,
			ImageURL:  image,
			SortOrder: i,
		}
		if err := global.DB.Create(taskImageModel).Error; err != nil {
			return err
		}
	}

	if err := global.DB.Create(taskInfoModel).Error; err != nil {
		return err
	}

	return nil
}

func GetTaskByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := global.DB.Preload("TaskInfo").Preload("TaskImages").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func GetTaskList(page, pageSize int) ([]models.TaskListResponse, error) {
	var tasks []models.TaskListResponse
	query := global.DB.Model(&models.Task{}).
		Joins("LEFT JOIN task_info ON tasks.id = task_info.task_id").
		Select("tasks.id as task_id, tasks.status, tasks.created_at, task_info.title, task_info.priority, task_info.due_date").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskAssigneeID(taskID uint) (uint, error) {
	var task models.Task
	if err := global.DB.Model(&models.Task{}).Where("id = ?", taskID).Select("assignee_id").First(&task).Error; err != nil {
		return 0, err
	}
	return task.AssigneeID, nil
}

func AssignTask(taskID, userID uint) error {
	if err := global.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("assignee_id", userID).Error; err != nil {
		return err
	}
	if err := global.DB.Create(&models.TaskHistory{
		TaskID:     taskID,
		Action:     "分配",
		FieldName:  "assignee_id",
		NewValue:   strconv.Itoa(int(userID)),
		OperatorID: userID,
		CreatedAt:  time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}
