import type { ApiResponse } from '../types/task'
import { http } from '../utils/http'

// 图片上传响应接口
export interface UploadImageResponse {
  path: string
}

// 支持的图片类型
export const SUPPORTED_IMAGE_TYPES = [
  'image/jpeg',
  'image/png', 
  'image/gif'
]

// 最大文件大小（5MB）
export const MAX_FILE_SIZE = 5 * 1024 * 1024

// 验证图片文件
export const validateImageFile = (file: File): string | null => {
  // 检查文件类型
  if (!SUPPORTED_IMAGE_TYPES.includes(file.type)) {
    return '只支持 JPEG、PNG、GIF 格式的图片'
  }
  
  // 检查文件大小
  if (file.size > MAX_FILE_SIZE) {
    return '图片大小不能超过 5MB'
  }
  
  return null
}

// 上传图片
export const uploadImage = async (file: File): Promise<ApiResponse<UploadImageResponse>> => {
  // 验证文件
  const validationError = validateImageFile(file)
  if (validationError) {
    throw new Error(validationError)
  }
  
  try {
    // 创建FormData
    const formData = new FormData()
    formData.append('file', file)
    
    // 发送请求
    const response = await fetch(`${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'}/img/`, {
      method: 'POST',
      headers: {
        'Authorization': `${localStorage.getItem('token')}`
      },
      body: formData
    })
    
    if (response.status === 401) {
      // 处理401错误，这会被http拦截器处理
      throw new Error('认证失败，请重新登录')
    }
    
    if (!response.ok) {
      throw new Error(`上传失败: ${response.status}`)
    }
    
    const result = await response.json()
    return result
  } catch (error) {
    console.error('图片上传失败:', error)
    throw error
  }
}

// 获取完整的图片URL
export const getImageUrl = (imagePath: string): string => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  // 如果路径已经是完整URL，直接返回
  if (imagePath.startsWith('http')) {
    return imagePath
  }
  // 否则拼接基础URL
  return `${baseUrl}${imagePath.startsWith('/') ? '' : '/'}${imagePath}`
} 