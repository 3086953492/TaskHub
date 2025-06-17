import type { ApiResponse } from '../types/task'

// API基础URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

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
  const token = localStorage.getItem('token')
  
  // 构建查询参数
  const searchParams = new URLSearchParams()
  searchParams.append('page', params.page.toString())
  searchParams.append('page_size', params.page_size.toString())
  
  if (params.assignee_id !== undefined) {
    searchParams.append('assignee_id', params.assignee_id.toString())
  }
  
  if (params.creator_id !== undefined) {
    searchParams.append('creator_id', params.creator_id.toString())
  }

  const response = await fetch(`${API_BASE_URL}/task?${searchParams.toString()}`, {
    method: 'GET',
    headers: {
      'Authorization': `${token}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return await response.json()
} 