import { useAuth } from '../composables/useAuth'
import router from '../router'
import { message } from './message'

// API基础URL
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// HTTP请求配置接口
export interface HttpConfig {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  data?: any
  params?: Record<string, string | number>
  headers?: Record<string, string>
  needAuth?: boolean
}

// HTTP响应接口
export interface HttpResponse<T = any> {
  data: T
  status: number
  statusText: string
}

// 创建HTTP客户端类
class HttpClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  // 处理401错误（token过期）
  private handleUnauthorized() {
    const { logout } = useAuth()
    
    // 清除认证信息
    logout()
    
    // 提示用户重新登录
    message.error('登录已过期', '您的登录状态已过期，请重新登录', 5000)
    
    // 延迟跳转，让用户看到提示消息
    setTimeout(() => {
      router.push('/login')
    }, 1000)
  }

  // 构建完整URL
  private buildURL(url: string, params?: Record<string, string | number>): string {
    const fullURL = url.startsWith('http') ? url : `${this.baseURL}${url}`
    
    if (params) {
      const searchParams = new URLSearchParams()
      Object.entries(params).forEach(([key, value]) => {
        searchParams.append(key, value.toString())
      })
      return `${fullURL}?${searchParams.toString()}`
    }
    
    return fullURL
  }

  // 构建请求头
  private buildHeaders(config: HttpConfig): Headers {
    const headers = new Headers()
    
    // 设置默认Content-Type
    headers.set('Content-Type', 'application/json')
    
    // 添加自定义头部
    if (config.headers) {
      Object.entries(config.headers).forEach(([key, value]) => {
        headers.set(key, value)
      })
    }
    
    // 添加认证头部
    if (config.needAuth !== false) {
      const { getAuthHeader } = useAuth()
      const authHeaders = getAuthHeader()
      Object.entries(authHeaders).forEach(([key, value]) => {
        headers.set(key, value)
      })
    }
    
    return headers
  }

  // 通用请求方法
  async request<T = any>(config: HttpConfig): Promise<HttpResponse<T>> {
    const { url, method = 'GET', data, params } = config
    
    const fullURL = this.buildURL(url, params)
    const headers = this.buildHeaders(config)
    
    const requestInit: RequestInit = {
      method,
      headers,
    }
    
    // 添加请求体
    if (data && ['POST', 'PUT', 'PATCH'].includes(method)) {
      requestInit.body = JSON.stringify(data)
    }
    
    try {
      const response = await fetch(fullURL, requestInit)
      
      // 检查是否是401错误
      if (response.status === 401) {
        this.handleUnauthorized()
        throw new Error('认证失败，请重新登录')
      }
      
      // 检查其他HTTP错误
      if (!response.ok) {
        let errorMessage = `请求失败: ${response.status}`
        try {
          const errorData = await response.json()
          if (errorData.message) {
            errorMessage = errorData.message
          }
        } catch (e) {
          // 如果无法解析错误响应，使用默认错误消息
        }
        throw new Error(errorMessage)
      }
      
      const responseData = await response.json()
      
      return {
        data: responseData,
        status: response.status,
        statusText: response.statusText
      }
    } catch (error) {
      // 如果是网络错误或其他错误，重新抛出
      if (error instanceof Error) {
        throw error
      }
      throw new Error('网络请求失败，请检查网络连接')
    }
  }

  // 便捷方法
  get<T = any>(url: string, config?: Omit<HttpConfig, 'url' | 'method'>): Promise<HttpResponse<T>> {
    return this.request<T>({ ...config, url, method: 'GET' })
  }

  post<T = any>(url: string, data?: any, config?: Omit<HttpConfig, 'url' | 'method' | 'data'>): Promise<HttpResponse<T>> {
    return this.request<T>({ ...config, url, method: 'POST', data })
  }

  put<T = any>(url: string, data?: any, config?: Omit<HttpConfig, 'url' | 'method' | 'data'>): Promise<HttpResponse<T>> {
    return this.request<T>({ ...config, url, method: 'PUT', data })
  }

  delete<T = any>(url: string, config?: Omit<HttpConfig, 'url' | 'method'>): Promise<HttpResponse<T>> {
    return this.request<T>({ ...config, url, method: 'DELETE' })
  }

  patch<T = any>(url: string, data?: any, config?: Omit<HttpConfig, 'url' | 'method' | 'data'>): Promise<HttpResponse<T>> {
    return this.request<T>({ ...config, url, method: 'PATCH', data })
  }
}

// 创建默认HTTP客户端实例
export const http = new HttpClient(API_BASE_URL)

// 导出类型
export { HttpClient } 