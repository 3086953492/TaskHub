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

func GetTaskDetail(id uint) (*models.TaskDetailResponse, error) {
	var task models.TaskDetailResponse
	if err := global.DB.Model(&models.Task{}).
		Joins("LEFT JOIN task_info ON tasks.id = task_info.task_id").
		Joins("LEFT JOIN task_images ON tasks.id = task_images.task_id").
		Select("tasks.id as task_id, tasks.status, tasks.created_at, tasks.updated_at, task_info.title, task_info.description, task_info.priority, task_info.due_date, task_info.completed_at, task_images.image_url").
		First(&task, id).Error; err != nil {
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

	if err := global.DB.Model(&models.Task{}).Where("id = ?", taskID).Update("assignee_id", userID).Update("status", 2).Error; err != nil {
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

func GetUnassignedTasks(page, pageSize int) ([]models.TaskListResponse, error) {

	var tasks []models.TaskListResponse

	query := global.DB.Model(&models.Task{}).
		Joins("LEFT JOIN task_info ON tasks.id = task_info.task_id").
		Select("tasks.id as task_id, tasks.status, tasks.created_at, task_info.title, task_info.priority, task_info.due_date").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Where("status = 1")

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// 获取任务历史记录基本信息
func GetTaskHistoryRecords(taskID uint) ([]models.TaskHistory, error) {
	var histories []models.TaskHistory
	if err := global.DB.Where("task_id = ?", taskID).Order("created_at DESC").Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

// 获取历史记录对应的图片
func GetTaskHistoryImages(historyID uint) ([]string, error) {
	var images []string
	var historyImages []models.TaskHistoryImage
	if err := global.DB.Where("history_id = ? AND deleted_at IS NULL", historyID).
		Order("sort_order ASC").Find(&historyImages).Error; err != nil {
		return nil, err
	}

	for _, img := range historyImages {
		images = append(images, img.ImageURL)
	}
	return images, nil
}
