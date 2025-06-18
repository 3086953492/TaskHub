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

        <!-- 头像上传区域 -->
        <div class="form-section">
          <h3>头像设置</h3>
          <p class="section-tip">支持 JPEG、PNG、GIF 格式，文件大小不超过 5MB</p>
          
          <div class="avatar-upload-area">
            <div class="avatar-preview">
              <img 
                v-if="getCurrentAvatarUrl()" 
                :src="getCurrentAvatarUrl()"
                :alt="user?.nickname || user?.username"
                class="avatar-image"
              />
              <div v-else class="avatar-placeholder">
                {{ getAvatarText(user?.nickname || user?.username || '') }}
              </div>
              
              <!-- 上传进度遮罩 -->
              <div v-if="uploadProgress !== null" class="upload-overlay">
                <div class="upload-progress">
                  <div class="progress-circle">
                    <svg class="progress-ring" width="60" height="60">
                      <circle
                        class="progress-ring-circle"
                        stroke="white"
                        stroke-width="4"
                        fill="transparent"
                        r="26"
                        cx="30"
                        cy="30"
                        :stroke-dasharray="163.36"
                        :stroke-dashoffset="163.36 - (163.36 * uploadProgress) / 100"
                      />
                    </svg>
                    <span class="progress-text">{{ uploadProgress }}%</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="upload-actions">
              <input
                ref="fileInput"
                type="file"
                accept="image/jpeg,image/png,image/gif"
                @change="handleFileSelect"
                style="display: none"
              />
              
              <button 
                type="button" 
                @click="triggerFileUpload" 
                class="upload-btn"
                :disabled="isLoading || uploadProgress !== null"
              >
                <span class="btn-icon">📁</span>
                <span>选择图片</span>
              </button>
              
              <button 
                v-if="showRemoveButton"
                type="button" 
                @click="removeAvatar" 
                class="remove-btn"
                :disabled="isLoading || uploadProgress !== null"
              >
                <span class="btn-icon">🗑️</span>
                <span>移除头像</span>
              </button>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { updateUserInfo } from '../api/auth'
import { uploadImage, getImageUrl, validateImageFile } from '../api/upload'
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
  email: '',
  avatar: ''
})

// 状态
const isLoading = ref(false)
const error = ref('')

// 头像相关状态
const avatarPreview = ref<string>('')
const uploadProgress = ref<number | null>(null)
const fileInput = ref<HTMLInputElement>()
const newAvatarPath = ref<string>('')

// 计算属性：是否有头像变化
const hasAvatarChange = computed(() => {
  return newAvatarPath.value !== '' || avatarPreview.value !== ''
})

// 计算属性：是否显示移除按钮
const showRemoveButton = computed(() => {
  // 有新上传的头像预览，或者用户原本有头像且没有标记移除
  return avatarPreview.value || (user.value?.avatar && newAvatarPath.value !== 'REMOVE_AVATAR')
})

// 获取头像文字
const getAvatarText = (name: string) => {
  return name.charAt(0).toUpperCase()
}

// 获取当前显示的头像URL
const getCurrentAvatarUrl = (): string => {
  // 如果标记移除头像，不显示任何图片
  if (newAvatarPath.value === 'REMOVE_AVATAR') {
    return ''
  }
  // 优先显示新上传的预览图
  if (avatarPreview.value) {
    return avatarPreview.value
  }
  // 如果用户原本有头像且没有被移除，显示原头像
  if (user.value?.avatar && !newAvatarPath.value) {
    return getImageUrl(user.value.avatar)
  }
  // 否则不显示图片
  return ''
}

// 头像相关方法
const triggerFileUpload = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return
  
  // 验证文件
  const validationError = validateImageFile(file)
  if (validationError) {
    error.value = validationError
    return
  }
  
  // 显示预览
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
  
  // 上传文件
  uploadAvatarFile(file)
}

const uploadAvatarFile = async (file: File) => {
  uploadProgress.value = 0
  error.value = ''
  
  try {
    // 模拟上传进度
    const progressInterval = setInterval(() => {
      if (uploadProgress.value !== null && uploadProgress.value < 90) {
        uploadProgress.value += 10
      }
    }, 100)
    
    const response = await uploadImage(file)
    
    clearInterval(progressInterval)
    uploadProgress.value = 100
    
    if (response.code === 200 && response.data) {
      newAvatarPath.value = response.data.path
      message.success('上传成功', '头像上传完成')
      
      // 延迟清除进度
      setTimeout(() => {
        uploadProgress.value = null
      }, 500)
    } else {
      throw new Error(response.msg || '上传失败')
    }
  } catch (err) {
    uploadProgress.value = null
    avatarPreview.value = ''
    newAvatarPath.value = ''
    
    const errorMessage = err instanceof Error ? err.message : '上传失败'
    error.value = errorMessage
    message.error('上传失败', errorMessage)
  }
}

const removeAvatar = () => {
  // 如果有新上传的头像，清除预览和新路径
  if (avatarPreview.value || newAvatarPath.value) {
    avatarPreview.value = ''
    newAvatarPath.value = ''
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  } else if (user.value?.avatar) {
    // 如果要移除原有头像，设置移除标记
    newAvatarPath.value = 'REMOVE_AVATAR'
  }
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
    originalData.avatar = user.value.avatar || ''
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

  // 检查头像是否有修改
  if (newAvatarPath.value === 'REMOVE_AVATAR') {
    // 如果标记移除头像，传空字符串
    changes.avatar = ''
  } else if (newAvatarPath.value && newAvatarPath.value !== 'REMOVE_AVATAR') {
    // 如果有新上传的头像，使用完整URL
    changes.avatar = getImageUrl(newAvatarPath.value)
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
  // 恢复头像状态
  resetAvatarState()
  router.back()
}

// 重置头像状态
const resetAvatarState = () => {
  avatarPreview.value = ''
  newAvatarPath.value = ''
  uploadProgress.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
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
  align-items: flex-start;
  gap: 20px;
}

.avatar-preview {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  position: relative;
  border: 3px solid #e5e7eb;
}

.avatar-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.upload-progress {
  position: relative;
}

.progress-circle {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.progress-ring {
  transform: rotate(-90deg);
}

.progress-ring-circle {
  transition: stroke-dashoffset 0.3s ease;
}

.progress-text {
  position: absolute;
  color: white;
  font-size: 12px;
  font-weight: 600;
}

.upload-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.upload-btn,
.remove-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.upload-btn {
  background: #3b82f6;
  color: white;
}

.upload-btn:hover:not(:disabled) {
  background: #2563eb;
}

.upload-btn:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.remove-btn {
  background: #fee2e2;
  color: #dc2626;
}

.remove-btn:hover:not(:disabled) {
  background: #fecaca;
}

.remove-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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