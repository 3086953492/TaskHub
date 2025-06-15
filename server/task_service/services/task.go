package services

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
	"TaskHub/task_service/repositories"
	"TaskHub/task_service/utils/logger"
	"errors"
	"fmt"
	"strconv"
)

func CreateTask(task *models.CreateTaskRequest, creatorID uint) error {
	// 参数验证
	if err := global.Validate.Struct(task); err != nil {
		return err
	}

	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建任务基本记录
	taskRecord, err := repositories.CreateTaskRecord(tx, creatorID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 创建任务信息记录
	if err := repositories.CreateTaskInfo(tx, taskRecord.ID, task); err != nil {
		tx.Rollback()
		return err
	}

	// 创建任务图片记录（如果有图片）
	if len(task.Images) > 0 {
		if err := repositories.CreateTaskImages(tx, taskRecord.ID, task.Images); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("创建任务成功, 任务ID: %d, 创建人ID: %d，任务: %v", taskRecord.ID, creatorID, task))
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

	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新任务分配人
	if err := repositories.UpdateTaskAssignee(tx, taskID, userID); err != nil {
		tx.Rollback()
		return err
	}

	// 更新任务状态为已分配(2)
	if err := repositories.UpdateTaskStatus(tx, taskID, 2); err != nil {
		tx.Rollback()
		return err
	}

	// 创建历史记录
	_, err = repositories.CreateTaskHistory(tx, taskID, "分配", strconv.Itoa(int(userID)), userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("任务分配成功, 任务ID: %d，分配人ID: %d", taskID, userID))

	return nil
}

func GetTaskList(page, pageSize uint) ([]models.TaskListResponse, error) {
	// 获取任务基本记录
	tasks, err := repositories.GetTaskRecords(int(page), int(pageSize), nil)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []models.TaskListResponse{}, nil
	}

	// 提取任务ID列表
	taskIDs := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	// 批量获取任务信息
	taskInfoMap, err := repositories.GetTaskInfoBatch(taskIDs)
	if err != nil {
		return nil, err
	}

	// 组装响应数据
	var result []models.TaskListResponse
	for _, task := range tasks {
		taskInfo := taskInfoMap[task.ID]
		if taskInfo == nil {
			// 如果没有找到对应的TaskInfo，跳过或使用默认值
			continue
		}

		response := models.TaskListResponse{
			TaskID:    task.ID,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
			Title:     taskInfo.Title,
			Priority:  taskInfo.Priority,
			DueDate:   taskInfo.DueDate,
		}
		result = append(result, response)
	}

	return result, nil
}

func GetTaskDetail(taskID, userID uint, role string) (*models.TaskDetailResponse, error) {

	if !CheckTaskViewPermission(taskID, userID, role) {
		return nil, errors.New("无权限查看任务详情")
	}

	// 获取任务基本信息
	taskRecord, err := repositories.GetTaskRecord(taskID)
	if err != nil {
		return nil, err
	}

	// 获取任务详细信息
	taskInfo, err := repositories.GetTaskInfo(taskID)
	if err != nil {
		return nil, err
	}

	// 获取任务图片
	images, err := repositories.GetTaskImages(taskID)
	if err != nil {
		// 图片获取失败时，继续处理但记录空数组
		images = []models.TaskImageResponse{}
	}

	// 构建响应对象
	response := &models.TaskDetailResponse{
		TaskID:      taskRecord.ID,
		Status:      taskRecord.Status,
		AssigneeID:  taskRecord.AssigneeID,
		CreatorID:   taskRecord.CreatorID,
		CreatedAt:   taskRecord.CreatedAt,
		UpdatedAt:   taskRecord.UpdatedAt,
		Title:       taskInfo.Title,
		Description: taskInfo.Description,
		Priority:    taskInfo.Priority,
		DueDate:     taskInfo.DueDate,
		CompletedAt: taskInfo.CompletedAt,
		Images:      images,
	}

	return response, nil
}

func GetUnassignedTasks(page, pageSize uint) ([]models.TaskListResponse, error) {
	// 设置条件：只获取状态为1（待处理）的任务
	conditions := map[string]interface{}{
		"status": 1,
	}

	// 获取未分配的任务基本记录
	tasks, err := repositories.GetTaskRecords(int(page), int(pageSize), conditions)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []models.TaskListResponse{}, nil
	}

	// 提取任务ID列表
	taskIDs := make([]uint, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}

	// 批量获取任务信息
	taskInfoMap, err := repositories.GetTaskInfoBatch(taskIDs)
	if err != nil {
		return nil, err
	}

	// 组装响应数据
	var result []models.TaskListResponse
	for _, task := range tasks {
		taskInfo := taskInfoMap[task.ID]
		if taskInfo == nil {
			// 如果没有找到对应的TaskInfo，跳过或使用默认值
			continue
		}

		response := models.TaskListResponse{
			TaskID:    task.ID,
			Status:    task.Status,
			CreatedAt: task.CreatedAt,
			Title:     taskInfo.Title,
			Priority:  taskInfo.Priority,
			DueDate:   taskInfo.DueDate,
		}
		result = append(result, response)
	}

	return result, nil
}

func GetTaskHistory(taskID, userID uint, role string) ([]models.TaskHistoryResponse, error) {

	if !CheckTaskViewPermission(taskID, userID, role) {
		return nil, errors.New("无权限查看任务历史记录")
	}

	// 获取历史记录基本信息
	histories, err := repositories.GetTaskHistoryRecords(taskID)
	if err != nil {
		return nil, err
	}

	var result []models.TaskHistoryResponse
	for _, history := range histories {
		// 获取该历史记录对应的图片
		images, err := repositories.GetTaskHistoryImages(history.ID)
		if err != nil {
			// 图片获取失败时，继续处理但记录空数组
			images = []string{}
		}

		// 构建响应对象
		response := models.TaskHistoryResponse{
			ID:         history.ID,
			Action:     history.Action,
			OperatorID: history.OperatorID,
			Remark:     history.Remark,
			CreatedAt:  history.CreatedAt,
			Images:     images,
		}
		result = append(result, response)
	}

	return result, nil
}

func CheckTaskViewPermission(taskID, userID uint, role string) bool {

	if role != "admin" {
		taskAssigneeID, err := repositories.GetTaskAssigneeID(taskID)
		if err != nil {
			return false
		}
		taskCreatorID, err := repositories.GetTaskCreatorID(taskID)
		if err != nil {
			return false
		}
		if taskAssigneeID != userID && taskCreatorID != userID && taskAssigneeID != 0 { // 如果任务已分配，则只有创建者和分配者可以查看
			return false
		}
	}

	return true
}

func UpdateTask(taskID, userID uint, role string, updateTaskRequest *models.UpdateTaskRequest) error {
	if !CheckTaskUpdatePermission(taskID, userID, role) {
		return errors.New("无权限更新任务")
	}

	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新任务状态
	if updateTaskRequest.Status != 0 {
		if err := repositories.UpdateTaskStatus(tx, taskID, updateTaskRequest.Status); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 更新任务信息
	if err := repositories.UpdateTaskInfo(tx, taskID, &updateTaskRequest.UpdateTaskInfo); err != nil {
		tx.Rollback()
		return err
	}

	// 更新现有图片
	for _, updateImage := range updateTaskRequest.UpdateImages {
		if err := repositories.UpdateTaskImage(tx, updateImage.ImageID, &updateImage); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 删除图片
	for _, deleteImage := range updateTaskRequest.DeleteImages {
		if err := repositories.DeleteTaskImage(tx, deleteImage.ImageID); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 新增图片
	for _, newImage := range updateTaskRequest.NewImages {
		if err := repositories.CreateTaskImage(tx, taskID, &newImage); err != nil {
			tx.Rollback()
			return err
		}
	}

	history, err := repositories.CreateTaskHistory(tx, taskID, "更新任务信息", updateTaskRequest.Remark, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 如果有备注图片，插入历史记录图片
	if len(updateTaskRequest.RemarkImages) > 0 {
		if err := repositories.CreateTaskHistoryImages(tx, history.ID, updateTaskRequest.RemarkImages); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CheckTaskUpdatePermission(taskID, userID uint, role string) bool {

	if role != "admin" {
		taskCreatorID, err := repositories.GetTaskCreatorID(taskID)
		if err != nil {
			return false
		}
		if taskCreatorID != userID {
			return false
		}
	}
	return true
}

func UpdateTaskStatus(taskID, userID uint, role string, updateTaskRequest *models.UpdateTaskStatusRequest) error {
	if !CheckTaskStatusUpdatePermission(taskID, userID, role) {
		return errors.New("无权限更新任务状态")
	}

	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := repositories.UpdateTaskStatus(tx, taskID, updateTaskRequest.Status); err != nil {
		tx.Rollback()
		return err
	}

	history, err := repositories.CreateTaskHistory(tx, taskID, "更新任务状态", updateTaskRequest.Remark, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(updateTaskRequest.RemarkImages) > 0 {
		if err := repositories.CreateTaskHistoryImages(tx, history.ID, updateTaskRequest.RemarkImages); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CheckTaskStatusUpdatePermission(taskID, userID uint, role string) bool {

	if role != "admin" {
		taskAssigneeID, _ := repositories.GetTaskAssigneeID(taskID)
		taskCreatorID, _ := repositories.GetTaskCreatorID(taskID)
		if taskAssigneeID != userID && taskCreatorID != userID {
			return false
		}
	}
	return true
}