package repositories

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// 创建任务基本记录
func CreateTaskRecord(tx *gorm.DB, creatorID uint) (*models.Task, error) {
	taskModel := &models.Task{
		Status:    1,
		CreatorID: creatorID,
	}

	if err := tx.Create(taskModel).Error; err != nil {
		return nil, err
	}

	return taskModel, nil
}

// 创建任务信息记录
func CreateTaskInfo(tx *gorm.DB, taskID uint, req *models.CreateTaskRequest) error {
	taskInfoModel := &models.TaskInfo{
		TaskID:      taskID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	return tx.Create(taskInfoModel).Error
}

// 创建任务图片记录
func CreateTaskImages(tx *gorm.DB, taskID uint, images []string) error {
	for i, imageURL := range images {
		taskImageModel := &models.TaskImage{
			TaskID:    taskID,
			ImageURL:  imageURL,
			SortOrder: i,
		}
		if err := tx.Create(taskImageModel).Error; err != nil {
			return err
		}
	}
	return nil
}

// 获取任务基本信息
func GetTaskRecord(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := global.DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// 获取任务详细信息
func GetTaskInfo(taskID uint) (*models.TaskInfo, error) {
	var taskInfo models.TaskInfo
	if err := global.DB.Where("task_id = ?", taskID).First(&taskInfo).Error; err != nil {
		return nil, err
	}
	return &taskInfo, nil
}

// 获取任务图片列表
func GetTaskImages(taskID uint) ([]string, error) {
	var images []string
	var taskImages []models.TaskImage
	if err := global.DB.Where("task_id = ? AND deleted_at IS NULL", taskID).
		Order("sort_order ASC").Find(&taskImages).Error; err != nil {
		return nil, err
	}

	for _, img := range taskImages {
		images = append(images, img.ImageURL)
	}
	return images, nil
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
