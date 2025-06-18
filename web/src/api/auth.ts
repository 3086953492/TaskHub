import type { ApiResponse } from '../types/task'
import { http } from '../utils/http'

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

// 注册请求参数
export interface RegisterParams {
  username: string
  email: string
  password: string
  nickname: string
}

// 登录响应数据
export interface LoginResponse {
  user: User
  token: string
}

// 注册响应数据
export interface RegisterResponse {
  user: Omit<User, 'avatar'>
}

// 用户登录
export const loginUser = async (params: LoginParams): Promise<ApiResponse<LoginResponse>> => {
  try {
    const response = await http.post<ApiResponse<LoginResponse>>('/user/login', params, {
      needAuth: false // 登录请求不需要token
    })
    return response.data
  } catch (error) {
    console.error('登录失败:', error)
    throw error
  }
}

// 用户注册
export const registerUser = async (params: RegisterParams): Promise<ApiResponse<RegisterResponse>> => {
  try {
    const response = await http.post<ApiResponse<RegisterResponse>>('/user/register', params, {
      needAuth: false // 注册请求不需要token
    })
    return response.data
  } catch (error) {
    console.error('注册失败:', error)
    throw error
  }
}

// 刷新令牌
export const refreshToken = async (): Promise<ApiResponse<{ token: string }>> => {
  try {
    const response = await http.post<ApiResponse<{ token: string }>>('/auth/refresh')
    return response.data
  } catch (error) {
    console.error('刷新token失败:', error)
    throw error
  }
}

// 获取用户信息
export const getUserInfo = async (): Promise<ApiResponse<User>> => {
  try {
    const response = await http.get<ApiResponse<User>>('/user/info')
    return response.data
  } catch (error) {
    console.error('获取用户信息失败:', error)
    throw error
  }
}

// 更新用户信息
export const updateUserInfo = async (params: Partial<User>): Promise<ApiResponse<User>> => {
  try {
    const response = await http.put<ApiResponse<User>>('/user/info', params)
    return response.data
  } catch (error) {
    console.error('更新用户信息失败:', error)
    throw error
  }
} 