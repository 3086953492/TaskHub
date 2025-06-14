package services

import (
	"TaskHub/task_service/global"
	"TaskHub/task_service/models"
	"TaskHub/task_service/repositories"
	"TaskHub/task_service/utils/logger"
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

func GetTaskDetail(taskID, userID uint, role string) (*models.TaskDetailResponse, error) {
	if role != "admin" {
		taskAssigneeID, err := repositories.GetTaskAssigneeID(taskID)
		if err != nil {
			return nil, err
		}
		if taskAssigneeID != userID {
			return nil, errors.New("无权限查看任务详情")
		}
	}
	return repositories.GetTaskDetail(taskID)
}

func GetUnassignedTasks(page, pageSize uint) ([]models.TaskListResponse, error) {
	return repositories.GetUnassignedTasks(int(page), int(pageSize))
}

func GetTaskHistory(taskID, userID uint, role string) ([]models.TaskHistoryResponse, error) {
	if role != "admin" {
		taskAssigneeID, err := repositories.GetTaskAssigneeID(taskID)
		if err != nil {
			return nil, err
		}
		if taskAssigneeID != userID {
			return nil, errors.New("无权限查看任务详情")
		}
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
			FieldName:  history.FieldName,
			OldValue:   history.OldValue,
			NewValue:   history.NewValue,
			OperatorID: history.OperatorID,
			Remark:     history.Remark,
			CreatedAt:  history.CreatedAt,
			Images:     images,
		}
		result = append(result, response)
	}

	return result, nil
}
