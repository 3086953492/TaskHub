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
}

// 任务图片接口
export interface TaskImage {
  url: string
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