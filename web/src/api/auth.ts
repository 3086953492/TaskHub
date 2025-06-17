import type { ApiResponse } from '../types/task'

// API基础URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// 用户信息接口
export interface User {
  id: number
  username: string
  email: string
  role: string
  nickname: string
  avatar?: string
}

// 登录请求参数
export interface LoginParams {
  username: string
  password: string
}

// 登录响应数据
export interface LoginResponse {
  user: User
  token: string
}

// 用户登录
export const loginUser = async (params: LoginParams): Promise<ApiResponse<LoginResponse>> => {
  const response = await fetch(`${API_BASE_URL}/user/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(params),
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return await response.json()
}

// 刷新令牌
export const refreshToken = async (): Promise<ApiResponse<{ token: string }>> => {
  const token = localStorage.getItem('token')
  
  const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  return await response.json()
} 