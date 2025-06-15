package models

import "time"

// 任务表
type Task struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Status     int        `gorm:"default:1" json:"status"` // 状态 1:待处理 2:进行中 3:已完成 4:已取消
	AssigneeID uint       `json:"assignee_id"`             // 被分配人ID
	CreatorID  uint       `json:"creator_id"`              // 创建人ID
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at"`
}

func (t *Task) TableName() string {
	return "tasks"
}

// 任务信息表
type TaskInfo struct {
	TaskID      uint       `gorm:"uniqueIndex" json:"task_id"`
	Title       string     `gorm:"size:200;not null" json:"title"` // 任务标题
	Description string     `json:"description"`                    // 任务描述
	Priority    int        `gorm:"default:3" json:"priority"`      // 优先级 1:高 2:中 3:低
	DueDate     *time.Time `json:"due_date"`                       // 截止日期
	CompletedAt *time.Time `json:"completed_at"`                   // 完成时间
}

func (t *TaskInfo) TableName() string {
	return "task_info"
}

// 任务历史记录表
type TaskHistory struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID     uint      `json:"task_id"`
	Action     string    `gorm:"size:50;not null" json:"action"` // 操作类型
	OperatorID uint      `json:"operator_id"`                    // 操作人ID
	Remark     string    `gorm:"size:500" json:"remark"`         // 备注
	CreatedAt  time.Time `json:"created_at"`
}

func (t *TaskHistory) TableName() string {
	return "task_histories"
}

// 任务图像表
type TaskImage struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    uint       `json:"task_id"`
	ImageURL  string     `gorm:"size:500;not null" json:"image_url" ` // 图像URL路径
	SortOrder int        `gorm:"default:0" json:"sort_order"`         // 排序顺序
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (t *TaskImage) TableName() string {
	return "task_images"
}

// 任务历史记录图像表
type TaskHistoryImage struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	HistoryID uint       `json:"history_id"`
	ImageURL  string     `gorm:"size:500;not null" json:"image_url"` // 图像URL路径
	SortOrder int        `gorm:"default:0" json:"sort_order"`        // 排序顺序
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
}

func (t *TaskHistoryImage) TableName() string {
	return "task_history_images"
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    int        `json:"priority" binding:"required,oneof=1 2 3"`
	DueDate     *time.Time `json:"due_date"`
	Images      []string   `json:"images" validate:"imagesURL"`
}

type TaskListResponse struct {
	TaskID    uint       `json:"task_id"`
	Status    int        `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	Title     string     `json:"title"`
	Priority  int        `json:"priority"`
	DueDate   *time.Time `json:"due_date"`
}

// 任务图片响应结构
type TaskImageResponse struct {
	ID       uint   `json:"id"`
	ImageURL string `json:"image_url"`
}

type TaskDetailResponse struct {
	TaskID      uint                `json:"task_id"`
	Status      int                 `json:"status"`
	AssigneeID  uint                `json:"assignee_id"` // 被分配人ID
	CreatorID   uint                `json:"creator_id"`  // 创建人ID
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Priority    int                 `json:"priority"`
	DueDate     *time.Time          `json:"due_date"`
	CompletedAt *time.Time          `json:"completed_at"`
	Images      []TaskImageResponse `json:"images"`
}

type TaskHistoryResponse struct {
	ID         uint      `json:"id"`
	Action     string    `json:"action"`
	FieldName  string    `json:"field_name"`
	OldValue   string    `json:"old_value"`
	NewValue   string    `json:"new_value"`
	OperatorID uint      `json:"operator_id"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
	Images     []string  `json:"images"`
}

type UpdateTaskRequest struct {
	Status         int               `json:"status" binding:"oneof=1 4"`
	UpdateTaskInfo UpdateTaskInfo    `json:"update_task_info"`
	UpdateImages   []UpdateTaskImage `json:"update_images"`
	DeleteImages   []DeleteTaskImage `json:"delete_images"`
	NewImages      []NewTaskImage    `json:"new_images"`
	Remark         string            `json:"remark"`
	RemarkImages   []string          `json:"remark_images"` // 备注图片，插入历史记录图片表
}

type UpdateTaskInfo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority" binding:"oneof=1 2 3"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskImage struct { // 在更新任务信息时使用这个结构体，支持更新图片和排序
	ImageID   uint   `json:"image_id"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

type DeleteTaskImage struct { // 在更新任务信息时使用这个结构体，支持删除图片
	ImageID uint `json:"image_id"`
}

type NewTaskImage struct { // 在更新任务信息时使用这个结构体，支持新增图片
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`
}

type UpdateTaskStatusRequest struct {
	Status       int      `json:"status" binding:"oneof=2 3 4"`
	Remark       string   `json:"remark"`
	RemarkImages []string `json:"remark_images"`
}
