package repositories

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
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
func GetTaskImages(taskID uint) ([]models.TaskImageResponse, error) {
	var images []models.TaskImageResponse
	var taskImages []models.TaskImage
	if err := global.DB.Where("task_id = ? AND deleted_at IS NULL", taskID).
		Order("sort_order ASC").Find(&taskImages).Error; err != nil {
		return nil, err
	}

	for _, img := range taskImages {
		images = append(images, models.TaskImageResponse{
			ID:       img.ID,
			ImageURL: img.ImageURL,
		})
	}
	return images, nil
}

// 获取任务基本记录列表
func GetTaskRecords(page, pageSize int, conditions map[string]interface{}) ([]models.Task, error) {
	var tasks []models.Task
	query := global.DB.Model(&models.Task{}).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC")

	// 应用条件过滤
	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// 批量获取任务信息
func GetTaskInfoBatch(taskIDs []uint) (map[uint]*models.TaskInfo, error) {
	var taskInfos []models.TaskInfo
	if err := global.DB.Where("task_id IN ?", taskIDs).Find(&taskInfos).Error; err != nil {
		return nil, err
	}

	// 转换为map便于查找
	infoMap := make(map[uint]*models.TaskInfo)
	for i := range taskInfos {
		infoMap[taskInfos[i].TaskID] = &taskInfos[i]
	}
	return infoMap, nil
}

func GetTaskAssigneeID(taskID uint) (uint, error) {
	var task models.Task
	if err := global.DB.Model(&models.Task{}).Where("id = ?", taskID).Select("assignee_id").First(&task).Error; err != nil {
		return 0, err
	}
	return task.AssigneeID, nil
}

// 更新任务分配人
func UpdateTaskAssignee(tx *gorm.DB, taskID, assigneeID uint) error {
	db := tx
	if db == nil {
		db = global.DB
	}
	if err := db.Model(&models.Task{}).Where("id = ?", taskID).Update("assignee_id", assigneeID).Error; err != nil {
		return err
	}
	return nil
}

// 更新任务状态
func UpdateTaskStatus(tx *gorm.DB, taskID uint, status int) error {
	db := tx
	if db == nil {
		db = global.DB
	}
	if err := db.Model(&models.Task{}).Where("id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// 创建任务历史记录
func CreateTaskHistory(tx *gorm.DB, taskID uint, action, remark string, operatorID uint) (*models.TaskHistory, error) {
	db := tx
	if db == nil {
		db = global.DB
	}
	history := &models.TaskHistory{
		TaskID:     taskID,
		Action:     action,
		OperatorID: operatorID,
		Remark:     remark,
		CreatedAt:  time.Now(),
	}
	if err := db.Create(history).Error; err != nil {
		return nil, err
	}
	return history, nil
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

func GetTaskCreatorID(taskID uint) (uint, error) {
	var task models.Task
	if err := global.DB.Model(&models.Task{}).Where("id = ?", taskID).Select("creator_id").First(&task).Error; err != nil {
		return 0, err
	}
	return task.CreatorID, nil
}