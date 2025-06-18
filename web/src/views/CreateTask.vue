<template>
  <div class="create-task-container">
    <div class="page-header">
      <h1>创建新任务</h1>
      <p>添加新的任务项目</p>
    </div>

    <div class="task-form-container">
      <form @submit.prevent="handleSubmit" class="task-form">
        <!-- 基本信息 -->
        <div class="form-section">
          <h3>基本信息</h3>
          
          <div class="form-group">
            <label for="title" class="required">任务标题</label>
            <input
              id="title"
              v-model="form.title"
              type="text"
              placeholder="请输入任务标题"
              :disabled="isLoading"
              required
            />
          </div>

          <div class="form-group">
            <label for="description">任务描述</label>
            <textarea
              id="description"
              v-model="form.description"
              placeholder="请输入任务描述（可选）"
              :disabled="isLoading"
              rows="4"
            ></textarea>
          </div>

          <div class="form-group">
            <label for="priority" class="required">优先级</label>
            <select
              id="priority"
              v-model="form.priority"
              :disabled="isLoading"
              required
            >
              <option value="">请选择优先级</option>
              <option 
                v-for="option in PRIORITY_OPTIONS" 
                :key="option.value" 
                :value="option.value"
              >
                {{ option.label }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label for="dueDate">截止日期</label>
            <input
              id="dueDate"
              v-model="form.due_date"
              type="datetime-local"
              :disabled="isLoading"
            />
          </div>
        </div>

        <!-- 图片上传 -->
        <div class="form-section">
          <h3>任务图片</h3>
          <p class="section-tip">可以上传与任务相关的图片（可选）</p>
          
          <div class="image-upload-area">
            <!-- 已上传的图片列表 -->
            <div v-if="uploadedImages.length > 0" class="uploaded-images">
              <div 
                v-for="(image, index) in uploadedImages" 
                :key="index"
                class="image-item"
              >
                <img :src="image.url" :alt="`任务图片 ${index + 1}`" class="image-preview" />
                <button 
                  type="button" 
                  @click="removeImage(index)"
                  class="remove-image-btn"
                  :disabled="isLoading"
                >
                  ×
                </button>
              </div>
            </div>

            <!-- 上传按钮 -->
            <div class="upload-control">
              <input
                ref="imageInput"
                type="file"
                accept="image/jpeg,image/png,image/gif"
                @change="handleImageSelect"
                style="display: none"
                multiple
              />
              
              <button 
                type="button" 
                @click="triggerImageUpload" 
                class="upload-images-btn"
                :disabled="isLoading || imageUploadProgress !== null"
              >
                <span class="btn-icon">📷</span>
                <span>添加图片</span>
              </button>
              
              <div v-if="imageUploadProgress !== null" class="upload-progress-info">
                上传中... {{ imageUploadProgress }}%
              </div>
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
          <button type="submit" class="btn-primary" :disabled="isLoading || !form.title.trim()">
            <span v-if="isLoading">创建中...</span>
            <span v-else>创建任务</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { createTask, PRIORITY_OPTIONS } from '../api/task'
import type { CreateTaskParams, TaskImage } from '../api/task'
import { uploadImage, getImageUrl } from '../api/upload'
import { message } from '../utils/message'

const router = useRouter()

// 表单数据
const form = reactive<CreateTaskParams>({
  title: '',
  description: '',
  priority: 0,
  due_date: '',
  images: []
})

// 状态
const isLoading = ref(false)
const error = ref('')

// 图片相关状态
const uploadedImages = ref<TaskImage[]>([])
const imageUploadProgress = ref<number | null>(null)
const imageInput = ref<HTMLInputElement>()

// 触发图片上传
const triggerImageUpload = () => {
  imageInput.value?.click()
}

// 处理图片选择
const handleImageSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (!files || files.length === 0) return
  
  // 上传所有选中的图片
  for (let i = 0; i < files.length; i++) {
    await uploadTaskImage(files[i])
  }
  
  // 清空input
  if (imageInput.value) {
    imageInput.value.value = ''
  }
}

// 上传单个图片
const uploadTaskImage = async (file: File) => {
  imageUploadProgress.value = 0
  
  try {
    // 模拟上传进度
    const progressInterval = setInterval(() => {
      if (imageUploadProgress.value !== null && imageUploadProgress.value < 90) {
        imageUploadProgress.value += 10
      }
    }, 100)
    
    const response = await uploadImage(file)
    
    clearInterval(progressInterval)
    imageUploadProgress.value = 100
    
    if (response.code === 200 && response.data) {
      const imageUrl = getImageUrl(response.data.path)
      uploadedImages.value.push({ url: imageUrl })
      
      setTimeout(() => {
        imageUploadProgress.value = null
      }, 500)
    } else {
      throw new Error(response.msg || '图片上传失败')
    }
  } catch (err) {
    imageUploadProgress.value = null
    const errorMessage = err instanceof Error ? err.message : '图片上传失败'
    message.error('上传失败', errorMessage)
  }
}

// 移除图片
const removeImage = (index: number) => {
  uploadedImages.value.splice(index, 1)
}

// 验证表单
const validateForm = (): boolean => {
  error.value = ''
  
  if (!form.title.trim()) {
    error.value = '请输入任务标题'
    return false
  }
  
  if (!form.priority) {
    error.value = '请选择优先级'
    return false
  }
  
  // 验证截止日期格式
  if (form.due_date) {
    const dueDate = new Date(form.due_date)
    const now = new Date()
    
    if (dueDate <= now) {
      error.value = '截止日期必须晚于当前时间'
      return false
    }
  }
  
  return true
}

// 格式化日期为ISO格式
const formatDateToISO = (dateString: string): string => {
  if (!dateString) return ''
  return new Date(dateString).toISOString()
}

// 提交表单
const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }
  
  isLoading.value = true
  error.value = ''
  
  try {
    // 准备提交数据
    const taskData: CreateTaskParams = {
      title: form.title.trim(),
      description: form.description?.trim() || undefined,
      priority: form.priority,
      due_date: form.due_date ? formatDateToISO(form.due_date) : undefined,
      images: uploadedImages.value.length > 0 ? uploadedImages.value : undefined
    }
    
    const response = await createTask(taskData)
    
    if (response.code === 200) {
      message.success('创建成功', '任务已成功创建')
      
      // 延迟跳转到首页
      setTimeout(() => {
        router.push('/')
      }, 1500)
    } else {
      error.value = response.msg || '创建任务失败'
    }
  } catch (err) {
    console.error('创建任务失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 取消创建
const handleCancel = () => {
  router.back()
}
</script>

<style scoped>
.create-task-container {
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

.task-form-container {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.task-form {
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

.form-group label.required::after {
  content: ' *';
  color: #dc2626;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group input:disabled,
.form-group textarea:disabled,
.form-group select:disabled {
  background: #f9fafb;
  color: #6b7280;
  cursor: not-allowed;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.image-upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  transition: border-color 0.2s;
}

.image-upload-area:hover {
  border-color: #3b82f6;
}

.uploaded-images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.image-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid #e5e7eb;
}

.image-preview {
  width: 100%;
  height: 120px;
  object-fit: cover;
  display: block;
}

.remove-image-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.remove-image-btn:hover {
  background: rgba(0, 0, 0, 0.9);
}

.upload-control {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.upload-images-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.upload-images-btn:hover:not(:disabled) {
  background: #2563eb;
}

.upload-images-btn:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.upload-progress-info {
  color: #3b82f6;
  font-size: 14px;
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

.btn-icon {
  font-size: 14px;
}

@media (max-width: 768px) {
  .create-task-container {
    padding: 16px;
  }
  
  .task-form {
    padding: 24px;
  }
  
  .uploaded-images {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
    gap: 12px;
  }
  
  .image-preview {
    height: 100px;
  }
  
  .form-actions {
    flex-direction: column;
  }
}
</style> 