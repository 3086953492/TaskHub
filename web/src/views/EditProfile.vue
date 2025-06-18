<template>
  <div class="edit-profile-container">
    <div class="page-header">
      <h1>编辑个人信息</h1>
      <p>修改您的用户信息和密码</p>
    </div>

    <div class="edit-form-container">
      <form @submit.prevent="handleSubmit" class="edit-form">
        <!-- 基本信息 -->
        <div class="form-section">
          <h3>基本信息</h3>
          
          <div class="form-group">
            <label for="username">用户名</label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              placeholder="请输入用户名"
              :disabled="isLoading"
            />
          </div>

          <div class="form-group">
            <label for="nickname">昵称</label>
            <input
              id="nickname"
              v-model="form.nickname"
              type="text"
              placeholder="请输入昵称"
              :disabled="isLoading"
            />
          </div>

          <div class="form-group">
            <label for="email">邮箱</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              placeholder="请输入邮箱地址"
              :disabled="isLoading"
            />
          </div>
        </div>

        <!-- 密码修改 -->
        <div class="form-section">
          <h3>修改密码</h3>
          <p class="section-tip">如果不需要修改密码，请留空</p>
          
          <div class="form-group">
            <label for="currentPassword">当前密码</label>
            <input
              id="currentPassword"
              v-model="form.currentPassword"
              type="password"
              placeholder="请输入当前密码"
              :disabled="isLoading"
            />
          </div>

          <div class="form-group">
            <label for="newPassword">新密码</label>
            <input
              id="newPassword"
              v-model="form.newPassword"
              type="password"
              placeholder="请输入新密码（至少6位）"
              :disabled="isLoading"
            />
          </div>

          <div class="form-group">
            <label for="confirmPassword">确认新密码</label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              :disabled="isLoading"
            />
          </div>
        </div>

        <!-- 头像上传区域（暂时禁用） -->
        <div class="form-section disabled">
          <h3>头像设置</h3>
          <p class="section-tip">头像修改功能暂未开放</p>
          
          <div class="avatar-upload-area">
            <div class="avatar-preview">
              <div class="avatar-placeholder">
                {{ getAvatarText(user?.nickname || user?.username || '') }}
              </div>
            </div>
            <div class="upload-info">
              <p>头像修改功能开发中...</p>
            </div>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="error-message">
          {{ error }}
        </div>

        <!-- 操作按钮 -->
        <div class="form-actions">
          <button type="button" @click="handleCancel" class="btn-secondary" :disabled="isLoading">
            取消
          </button>
          <button type="submit" class="btn-primary" :disabled="isLoading">
            <span v-if="isLoading">保存中...</span>
            <span v-else>保存修改</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { updateUserInfo } from '../api/auth'
import { message } from '../utils/message'

const router = useRouter()
const { user, updateUser } = useAuth()

// 表单数据
const form = reactive({
  username: '',
  nickname: '',
  email: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 原始数据（用于比较是否有修改）
const originalData = reactive({
  username: '',
  nickname: '',
  email: ''
})

// 状态
const isLoading = ref(false)
const error = ref('')

// 获取头像文字
const getAvatarText = (name: string) => {
  return name.charAt(0).toUpperCase()
}

// 初始化表单数据
const initForm = () => {
  if (user.value) {
    form.username = user.value.username || ''
    form.nickname = user.value.nickname || ''
    form.email = user.value.email || ''
    
    // 保存原始数据
    originalData.username = user.value.username || ''
    originalData.nickname = user.value.nickname || ''
    originalData.email = user.value.email || ''
  }
}

// 验证表单
const validateForm = () => {
  error.value = ''

  // 检查是否有修改密码的操作
  const hasPasswordChange = form.currentPassword || form.newPassword || form.confirmPassword

  if (hasPasswordChange) {
    if (!form.currentPassword) {
      error.value = '请输入当前密码'
      return false
    }
    
    if (!form.newPassword) {
      error.value = '请输入新密码'
      return false
    }
    
    if (form.newPassword.length < 6) {
      error.value = '新密码至少需要6位'
      return false
    }
    
    if (form.newPassword !== form.confirmPassword) {
      error.value = '两次输入的新密码不一致'
      return false
    }
  }

  // 检查邮箱格式
  if (form.email && !isValidEmail(form.email)) {
    error.value = '请输入有效的邮箱地址'
    return false
  }

  return true
}

// 邮箱验证
const isValidEmail = (email: string) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// 获取修改的字段
const getChangedFields = () => {
  const changes: any = {}

  // 检查基本信息是否有修改
  if (form.username !== originalData.username) {
    changes.username = form.username
  }
  if (form.nickname !== originalData.nickname) {
    changes.nickname = form.nickname
  }
  if (form.email !== originalData.email) {
    changes.email = form.email
  }

  // 检查是否要修改密码
  if (form.newPassword) {
    changes.password = form.newPassword
  }

  return changes
}

// 提交表单
const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  const changes = getChangedFields()

  // 如果没有任何修改
  if (Object.keys(changes).length === 0) {
    message.info('没有修改', '没有检测到任何修改内容')
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    const response = await updateUserInfo(changes)
    
    if (response.code === 200) {
      // 更新本地用户信息（排除密码）
      const userUpdates = { ...changes }
      delete userUpdates.password
      
      if (Object.keys(userUpdates).length > 0) {
        updateUser(userUpdates)
      }
      
      message.success('保存成功', '用户信息已更新')
      
      // 延迟返回上一页
      setTimeout(() => {
        router.back()
      }, 1500)
    } else {
      error.value = response.msg || '保存失败'
    }
  } catch (err) {
    console.error('更新用户信息失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 取消编辑
const handleCancel = () => {
  router.back()
}

// 生命周期
onMounted(() => {
  initForm()
})
</script>

<style scoped>
.edit-profile-container {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background: #f8fafc;
  min-height: calc(100vh - 64px);
}

.page-header {
  text-align: center;
  margin-bottom: 32px;
}

.page-header h1 {
  color: #1f2937;
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.page-header p {
  color: #6b7280;
  font-size: 16px;
  margin: 0;
}

.edit-form-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.edit-form {
  padding: 32px;
}

.form-section {
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.form-section:last-of-type {
  border-bottom: none;
}

.form-section.disabled {
  opacity: 0.6;
  pointer-events: none;
}

.form-section h3 {
  color: #1f2937;
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 16px 0;
}

.section-tip {
  color: #6b7280;
  font-size: 14px;
  margin: -8px 0 16px 0;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  color: #374151;
  font-weight: 500;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group input:disabled {
  background: #f9fafb;
  color: #6b7280;
  cursor: not-allowed;
}

.avatar-upload-area {
  display: flex;
  align-items: center;
  gap: 20px;
}

.avatar-preview {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: 700;
}

.upload-info p {
  color: #6b7280;
  margin: 0;
}

.error-message {
  background: #fee2e2;
  color: #dc2626;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 20px;
  font-size: 14px;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 32px;
}

.btn-primary,
.btn-secondary {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 16px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-primary:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover:not(:disabled) {
  background: #e5e7eb;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .edit-profile-container {
    padding: 16px;
  }
  
  .edit-form {
    padding: 24px;
  }
  
  .avatar-upload-area {
    flex-direction: column;
    text-align: center;
  }
  
  .form-actions {
    flex-direction: column;
  }
}
</style> 