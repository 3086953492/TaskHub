import { ref, computed } from 'vue'
import type { User } from '../api/auth'

// 全局用户状态
const user = ref<User | null>(null)
const token = ref<string | null>(null)

// 初始化时从localStorage读取
const initAuth = () => {
  const savedToken = localStorage.getItem('token')
  const savedUser = localStorage.getItem('user')
  
  if (savedToken) {
    token.value = savedToken
  }
  
  if (savedUser) {
    try {
      user.value = JSON.parse(savedUser)
    } catch (error) {
      console.error('解析用户信息失败:', error)
      localStorage.removeItem('user')
    }
  }
}

export const useAuth = () => {
  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // 登录
  const login = (userData: User, userToken: string) => {
    user.value = userData
    token.value = userToken
    
    // 存储到localStorage
    localStorage.setItem('token', userToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  // 登出
  const logout = () => {
    user.value = null
    token.value = null
    
    // 清除localStorage
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 更新用户信息
  const updateUser = (userData: Partial<User>) => {
    if (user.value) {
      user.value = { ...user.value, ...userData }
      localStorage.setItem('user', JSON.stringify(user.value))
    }
  }

  // 获取认证头
  const getAuthHeader = () => {
    return token.value ? { Authorization: `${token.value}` } : {}
  }

  return {
    // 状态
    user: computed(() => user.value),
    token: computed(() => token.value),
    isLoggedIn,
    isAdmin,
    
    // 方法
    login,
    logout,
    updateUser,
    getAuthHeader,
    initAuth
  }
} 