import type { ApiResponse } from '../types/task'
import { http } from '../utils/http'

// 任务列表查询参数
export interface TaskListParams {
  page: number
  page_size: number
  assignee_id?: number
  creator_id?: number
}

// 任务列表项接口
export interface TaskListItem {
  task_id: number
  status: number
  created_at: string
  title: string
  priority: number
  due_date?: string
  assignee_id?: number
}

// 任务图片接口
export interface TaskImage {
  id?: number
  url: string
  sort_order?: number
}

// 创建任务请求参数
export interface CreateTaskParams {
  title: string
  description?: string
  priority: number
  due_date?: string
  images?: TaskImage[]
}

// 创建任务响应数据
export interface CreateTaskResponse {
  task_id: number
  title: string
  description?: string
  priority: number
  due_date?: string
  status: number
  created_at: string
}

// 任务详情接口
export interface TaskDetail {
  task_id: number
  status: number
  assignee_id?: number
  creator_id: number
  created_at: string
  updated_at: string
  title: string
  description?: string
  priority: number
  due_date?: string
  completed_at?: string
  images?: TaskImage[]
}

// 任务状态选项
export const STATUS_OPTIONS = [
  { value: 1, label: '待处理', color: '#6b7280' },
  { value: 2, label: '进行中', color: '#3b82f6' },
  { value: 3, label: '已完成', color: '#10b981' },
  { value: 4, label: '已取消', color: '#ef4444' }
]

// 任务历史记录接口
export interface TaskHistory {
  id: number
  action: string
  operator_id: number
  remark: string
  created_at: string
  images: string[]
}

// 优先级选项
export const PRIORITY_OPTIONS = [
  { value: 1, label: '高优先级', color: '#dc2626' },
  { value: 2, label: '中优先级', color: '#f59e0b' },
  { value: 3, label: '低优先级', color: '#10b981' }
]

// 获取任务列表
export const getTaskList = async (params: TaskListParams): Promise<ApiResponse<TaskListItem[]>> => {
  try {
    const response = await http.get<ApiResponse<TaskListItem[]>>('/task', {
      params: {
        page: params.page,
        page_size: params.page_size,
        ...(params.assignee_id !== undefined && { assignee_id: params.assignee_id }),
        ...(params.creator_id !== undefined && { creator_id: params.creator_id })
      }
    })
    return response.data
  } catch (error) {
    console.error('获取任务列表失败:', error)
    throw error
  }
}

// 创建任务
export const createTask = async (params: CreateTaskParams): Promise<ApiResponse<CreateTaskResponse>> => {
  try {
    const response = await http.post<ApiResponse<CreateTaskResponse>>('/task', params)
    return response.data
  } catch (error) {
    console.error('创建任务失败:', error)
    throw error
  }
}

// 获取任务详情
export const getTaskDetail = async (taskId: number): Promise<ApiResponse<TaskDetail>> => {
  try {
    const response = await http.get<ApiResponse<TaskDetail>>(`/task/${taskId}`)
    return response.data
  } catch (error) {
    console.error('获取任务详情失败:', error)
    throw error
  }
}

// 工具函数：根据优先级获取样式
export const getPriorityOption = (priority: number) => {
  return PRIORITY_OPTIONS.find(option => option.value === priority) || PRIORITY_OPTIONS[1]
}

// 获取任务历史记录
export const getTaskHistory = async (taskId: number): Promise<ApiResponse<TaskHistory[]>> => {
  try {
    const response = await http.get<ApiResponse<TaskHistory[]>>(`/task/history/${taskId}`)
    return response.data
  } catch (error) {
    console.error('获取任务历史记录失败:', error)
    throw error
  }
}

// 分配任务给当前用户
export const assignTask = async (taskId: number): Promise<ApiResponse<null>> => {
  try {
    const response = await http.patch<ApiResponse<null>>(`/task?id=${taskId}`)
    return response.data
  } catch (error) {
    console.error('分配任务失败:', error)
    throw error
  }
}

// 工具函数：根据状态获取样式
export const getStatusOption = (status: number) => {
  return STATUS_OPTIONS.find(option => option.value === status) || STATUS_OPTIONS[0]
}

// 更新任务参数接口
export interface UpdateTaskParams {
  status?: number
  update_task_info?: {
    title?: string
    description?: string
    priority?: number
    due_date?: string
  }
  update_images?: Array<{
    id: number
    url: string
    sort_order: number
  }>
  delete_images?: Array<{
    id: number
  }>
  add_images?: Array<{
    url: string
    sort_order: number
  }>
  remark?: string
  remark_images?: Array<{
    url: string
  }>
}

// 更新任务
export const updateTask = async (taskId: number, params: UpdateTaskParams): Promise<ApiResponse<null>> => {
  try {
    const response = await http.patch<ApiResponse<null>>(`/task/${taskId}`, params)
    return response.data
  } catch (error) {
    console.error('更新任务失败:', error)
    throw error
  }
}

// 状态变更参数接口
export interface UpdateTaskStatusParams {
  status: number
  remark?: string
  remark_images?: Array<{
    url: string
  }>
}

// 更新任务状态
export const updateTaskStatus = async (taskId: number, params: UpdateTaskStatusParams): Promise<ApiResponse<null>> => {
  try {
    const response = await http.patch<ApiResponse<null>>(`/task/status/${taskId}`, params)
    return response.data
  } catch (error) {
    console.error('更新任务状态失败:', error)
    throw error
  }
} 