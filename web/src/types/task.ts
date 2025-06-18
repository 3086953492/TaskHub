// 任务接口定义
export interface Task {
  task_id: number
  status: TaskStatus
  created_at: string
  title: string
  priority: TaskPriority
  due_date?: string
}

// 任务详情接口定义（包含更多字段）
export interface TaskDetail extends Task {
  assignee_id?: number
  creator_id: number
  updated_at: string
  description?: string
  completed_at?: string
  images?: TaskImage[]
}

// 任务图片接口定义
export interface TaskImage {
  id?: number
  url: string
  sort_order?: number
}

// 任务历史记录接口定义
export interface TaskHistory {
  id: number
  action: string
  operator_id: number
  remark?: string
  created_at: string
  images?: string[]
}

// 任务状态枚举
export enum TaskStatus {
  PENDING = 1,      // 待处理
  IN_PROGRESS = 2,  // 进行中
  COMPLETED = 3,    // 已完成
  CANCELLED = 4     // 已取消
}

// 任务优先级枚举
export enum TaskPriority {
  HIGH = 1,     // 高优先级
  MEDIUM = 2,   // 中优先级
  LOW = 3       // 低优先级
}

// API响应接口定义
export interface ApiResponse<T = any> {
  code: number
  msg: string
  data?: T
  total_pages?: number
}

// 任务列表查询参数
export interface TaskListParams {
  page: number
  page_size: number
  assignee_id?: number
  creator_id?: number
}

// 创建任务参数
export interface CreateTaskParams {
  title: string
  description?: string
  priority: TaskPriority
  due_date?: string
  images?: TaskImage[]
}

// 更新任务参数
export interface UpdateTaskParams {
  status?: TaskStatus
  update_task_info?: {
    title?: string
    description?: string
    priority?: TaskPriority
    due_date?: string
  }
  update_images?: TaskImage[]
  delete_images?: { id: number }[]
  add_images?: TaskImage[]
  remark?: string
  remark_images?: TaskImage[]
}

// 任务状态文本映射
export const TASK_STATUS_TEXT = {
  [TaskStatus.PENDING]: '待处理',
  [TaskStatus.IN_PROGRESS]: '进行中',
  [TaskStatus.COMPLETED]: '已完成',
  [TaskStatus.CANCELLED]: '已取消'
} as const

// 任务优先级文本映射
export const TASK_PRIORITY_TEXT = {
  [TaskPriority.HIGH]: '高优先级',
  [TaskPriority.MEDIUM]: '中优先级',
  [TaskPriority.LOW]: '低优先级'
} as const 