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