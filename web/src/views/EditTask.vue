<template>
  <div class="edit-task-container">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载任务信息...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">❌</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button @click="fetchTaskDetail" class="retry-btn">重试</button>
    </div>

    <!-- 权限不足 -->
    <div v-else-if="!hasEditPermission" class="permission-denied">
      <div class="permission-icon">🔒</div>
      <h3>权限不足</h3>
      <p>只有任务创建者或管理员可以编辑任务</p>
      <button @click="goBack" class="back-btn">返回</button>
    </div>

    <!-- 编辑表单 -->
    <div v-else-if="taskDetail" class="edit-form-container">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-actions">
          <button @click="goBack" class="back-btn">
            <span class="btn-icon">←</span>
            <span>返回</span>
          </button>
        </div>
        
        <div class="header-title">
          <h1>编辑任务 #{{ taskDetail.task_id }}</h1>
          <p class="header-subtitle">{{ taskDetail.title }}</p>
        </div>
      </div>

      <!-- 编辑表单 -->
      <form @submit.prevent="handleSubmit" class="edit-form">
        <!-- 基本信息编辑 -->
        <div class="form-section">
          <h3 class="section-title">
            <span class="section-icon">📝</span>
            基本信息
          </h3>
          
          <div class="form-grid">
            <div class="form-group">
              <label for="title">任务标题</label>
              <input
                id="title"
                v-model="formData.update_task_info.title"
                type="text"
                placeholder="请输入任务标题"
                class="form-input"
              />
            </div>

            <div class="form-group">
              <label for="priority">优先级</label>
              <select id="priority" v-model="formData.update_task_info.priority" class="form-select">
                <option :value="undefined">保持原有优先级</option>
                <option v-for="option in PRIORITY_OPTIONS" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
            </div>

            <div class="form-group">
              <label for="due_date">截止日期</label>
              <input
                id="due_date"
                v-model="formData.update_task_info.due_date"
                type="datetime-local"
                class="form-input"
              />
            </div>

            <div class="form-group">
              <label for="status">任务状态</label>
              <select id="status" v-model="formData.status" class="form-select">
                <option :value="undefined">保持原有状态</option>
                <option :value="1">待处理</option>
                <option :value="4">已取消</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label for="description">任务描述</label>
            <textarea
              id="description"
              v-model="formData.update_task_info.description"
              rows="4"
              placeholder="请输入任务描述"
              class="form-textarea"
            ></textarea>
          </div>
        </div>

        <!-- 图片管理 -->
        <div class="form-section">
          <h3 class="section-title">
            <span class="section-icon">🖼️</span>
            图片管理
          </h3>

          <!-- 现有图片 -->
          <div v-if="taskDetail.images && taskDetail.images.length > 0" class="existing-images">
            <h4>现有图片</h4>
            <div class="images-grid">
              <div 
                v-for="(image, index) in taskDetail.images" 
                :key="image.id || index"
                class="image-item"
              >
                <img :src="image.url" :alt="`图片 ${index + 1}`" class="image-preview" />
                <div class="image-actions">
                  <button 
                    type="button" 
                    @click="startUpdateImage(image, index)"
                    class="action-btn update-btn"
                  >
                    更新
                  </button>
                  <button 
                    type="button" 
                    @click="deleteImage(image.id, index)"
                    class="action-btn delete-btn"
                  >
                    删除
                  </button>
                </div>
                <div class="image-order">排序: {{ image.sort_order || index + 1 }}</div>
              </div>
            </div>
          </div>

          <!-- 图片更新表单 -->
          <div v-if="updatingImageIndex !== null" class="update-image-form">
            <h4>更新图片</h4>
            <div class="form-group">
              <label>新图片URL</label>
              <input
                v-model="updateImageForm.url"
                type="url"
                placeholder="请输入新图片URL"
                class="form-input"
              />
            </div>
            <div class="form-group">
              <label>排序值</label>
              <input
                v-model="updateImageForm.sort_order"
                type="number"
                min="1"
                placeholder="排序值"
                class="form-input"
              />
            </div>
            <div class="form-actions">
              <button type="button" @click="confirmUpdateImage" class="confirm-btn">确认更新</button>
              <button type="button" @click="cancelUpdateImage" class="cancel-btn">取消</button>
            </div>
          </div>

          <!-- 添加新图片 -->
          <div class="add-images-section">
            <h4>添加新图片</h4>
            <div v-for="(newImage, index) in formData.add_images" :key="index" class="new-image-item">
              <div class="form-group">
                <label>图片URL</label>
                <input
                  v-model="newImage.url"
                  type="url"
                  placeholder="请输入图片URL"
                  class="form-input"
                />
              </div>
              <div class="form-group">
                <label>排序值</label>
                <input
                  v-model="newImage.sort_order"
                  type="number"
                  min="1"
                  placeholder="排序值"
                  class="form-input"
                />
              </div>
              <button type="button" @click="removeNewImage(index)" class="remove-btn">移除</button>
            </div>
            <button type="button" @click="addNewImage" class="add-btn">+ 添加图片</button>
          </div>
        </div>

        <!-- 备注信息 -->
        <div class="form-section">
          <h3 class="section-title">
            <span class="section-icon">💬</span>
            备注信息
          </h3>
          
          <div class="form-group">
            <label for="remark">备注文字</label>
            <textarea
              id="remark"
              v-model="formData.remark"
              rows="3"
              placeholder="请输入备注信息"
              class="form-textarea"
            ></textarea>
          </div>

          <!-- 备注图片 -->
          <div class="remark-images-section">
            <h4>备注图片</h4>
            <div v-for="(remarkImage, index) in formData.remark_images" :key="index" class="remark-image-item">
              <div class="form-group">
                <label>图片URL</label>
                <input
                  v-model="remarkImage.url"
                  type="url"
                  placeholder="请输入备注图片URL"
                  class="form-input"
                />
              </div>
              <button type="button" @click="removeRemarkImage(index)" class="remove-btn">移除</button>
            </div>
            <button type="button" @click="addRemarkImage" class="add-btn">+ 添加备注图片</button>
          </div>
        </div>

        <!-- 提交按钮 -->
        <div class="form-actions">
          <button type="submit" :disabled="isSubmitting" class="submit-btn">
            <span v-if="isSubmitting">提交中...</span>
            <span v-else>保存更改</span>
          </button>
          <button type="button" @click="resetForm" class="reset-btn">重置</button>
        </div>
      </form>
    </div>

    <!-- 成功消息 -->
    <div v-if="showSuccessMessage" class="success-message">
      <div class="success-icon">✅</div>
      <p>任务更新成功！</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getTaskDetail, updateTask, PRIORITY_OPTIONS, type TaskDetail, type UpdateTaskParams } from '../api/task'
import { useAuth } from '../composables/useAuth'

const route = useRoute()
const router = useRouter()
const { user, isAdmin } = useAuth()

// 响应式状态
const isLoading = ref(true)
const error = ref('')
const taskDetail = ref<TaskDetail | null>(null)
const isSubmitting = ref(false)
const showSuccessMessage = ref(false)

// 图片更新相关
const updatingImageIndex = ref<number | null>(null)
const updateImageForm = reactive({
  url: '',
  sort_order: 1
})

// 表单数据
const formData = reactive({
  update_task_info: {
    title: '',
    description: '',
    priority: undefined as number | undefined,
    due_date: ''
  },
  update_images: [] as Array<{id: number, url: string, sort_order: number}>,
  delete_images: [] as Array<{id: number}>,
  add_images: [] as Array<{url: string, sort_order: number}>,
  status: undefined as number | undefined,
  remark: '',
  remark_images: [] as Array<{url: string}>
})

// 权限检查
const hasEditPermission = computed(() => {
  if (!taskDetail.value || !user.value) return false
  return isAdmin.value || taskDetail.value.creator_id === user.value.id
})

// 获取任务详情
const fetchTaskDetail = async () => {
  const taskId = Number(route.params.id)
  if (!taskId) {
    error.value = '无效的任务ID'
    isLoading.value = false
    return
  }

  try {
    isLoading.value = true
    error.value = ''
    
    const response = await getTaskDetail(taskId)
    if (response.code === 200 && response.data) {
      taskDetail.value = response.data
      initFormData()
    } else {
      error.value = response.msg || '获取任务详情失败'
    }
  } catch (err: any) {
    error.value = err.message || '网络错误，请重试'
  } finally {
    isLoading.value = false
  }
}

// 初始化表单数据
const initFormData = () => {
  if (!taskDetail.value) return
  
  if (formData.update_task_info) {
    formData.update_task_info.title = taskDetail.value.title
    formData.update_task_info.description = taskDetail.value.description || ''
    formData.update_task_info.priority = undefined // 不预设优先级，让用户选择
    formData.update_task_info.due_date = taskDetail.value.due_date 
      ? new Date(taskDetail.value.due_date).toISOString().slice(0, 16) 
      : ''
  }
  // 确保状态也是 undefined
  formData.status = undefined
}

// 重置表单
const resetForm = () => {
  initFormData()
  if (formData.update_images) formData.update_images = []
  if (formData.delete_images) formData.delete_images = []
  if (formData.add_images) formData.add_images = []
  formData.remark = ''
  if (formData.remark_images) formData.remark_images = []
}

// 添加新图片
const addNewImage = () => {
  if (!formData.add_images) formData.add_images = []
  formData.add_images.push({
    url: '',
    sort_order: formData.add_images.length + 1
  })
}

// 移除新图片
const removeNewImage = (index: number) => {
  if (formData.add_images) {
    formData.add_images.splice(index, 1)
  }
}

// 添加备注图片
const addRemarkImage = () => {
  if (!formData.remark_images) formData.remark_images = []
  formData.remark_images.push({ url: '' })
}

// 移除备注图片
const removeRemarkImage = (index: number) => {
  if (formData.remark_images) {
    formData.remark_images.splice(index, 1)
  }
}

// 开始更新图片
const startUpdateImage = (image: any, index: number) => {
  updatingImageIndex.value = index
  updateImageForm.url = image.url
  updateImageForm.sort_order = image.sort_order || index + 1
}

// 确认更新图片
const confirmUpdateImage = () => {
  if (updatingImageIndex.value === null) return
  
  const currentImage = taskDetail.value?.images?.[updatingImageIndex.value]
  if (!currentImage?.id) return

  if (!formData.update_images) formData.update_images = []
  formData.update_images.push({
    id: currentImage.id,
    url: updateImageForm.url,
    sort_order: updateImageForm.sort_order
  })

  cancelUpdateImage()
}

// 取消更新图片
const cancelUpdateImage = () => {
  updatingImageIndex.value = null
  updateImageForm.url = ''
  updateImageForm.sort_order = 1
}

// 删除图片
const deleteImage = (imageId: number | undefined, index: number) => {
  if (!imageId) return
  
  if (confirm('确定要删除这张图片吗？')) {
    if (!formData.delete_images) formData.delete_images = []
    formData.delete_images.push({ id: imageId })
    
    // 如果当前正在更新这张图片，取消更新
    if (updatingImageIndex.value === index) {
      cancelUpdateImage()
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!taskDetail.value) return

  try {
    isSubmitting.value = true
    
    // 构建提交数据，只包含有值的字段
    const submitData: UpdateTaskParams = {}
    
    // 基本信息更新
    const updateInfo: any = {}
    if (formData.update_task_info?.title && formData.update_task_info.title.trim() !== taskDetail.value.title) {
      updateInfo.title = formData.update_task_info.title.trim()
    }
    if (formData.update_task_info?.description !== (taskDetail.value.description || '')) {
      updateInfo.description = formData.update_task_info.description
    }
    // 只有明确选择了不同的优先级才发送
    if (formData.update_task_info?.priority && 
        typeof formData.update_task_info.priority === 'number' &&
        [1, 2, 3].includes(formData.update_task_info.priority) && 
        formData.update_task_info.priority !== taskDetail.value.priority) {
      updateInfo.priority = formData.update_task_info.priority
    }
    if (formData.update_task_info?.due_date && formData.update_task_info.due_date.trim() !== '') {
      updateInfo.due_date = new Date(formData.update_task_info.due_date).toISOString()
    }
    
    if (Object.keys(updateInfo).length > 0) {
      submitData.update_task_info = updateInfo
    }

    // 状态更新 (只有明确选择了有效状态才发送)
    if (formData.status !== undefined && 
        typeof formData.status === 'number' &&
        (formData.status === 1 || formData.status === 4) &&
        formData.status !== taskDetail.value.status) {
      submitData.status = formData.status
    }

    // 图片更新
    if (formData.update_images && formData.update_images.length > 0) {
      submitData.update_images = formData.update_images
    }
    if (formData.delete_images && formData.delete_images.length > 0) {
      submitData.delete_images = formData.delete_images
    }
    if (formData.add_images && formData.add_images.some(img => img.url)) {
      submitData.add_images = formData.add_images.filter(img => img.url)
    }

    // 备注信息
    if (formData.remark) {
      submitData.remark = formData.remark
    }
    if (formData.remark_images && formData.remark_images.some(img => img.url)) {
      submitData.remark_images = formData.remark_images.filter(img => img.url)
    }

    // 检查是否有更新内容
    if (Object.keys(submitData).length === 0) {
      alert('没有需要更新的内容')
      return
    }

    const response = await updateTask(taskDetail.value.task_id, submitData)
    
    if (response.code === 200) {
      showSuccessMessage.value = true
      setTimeout(() => {
        showSuccessMessage.value = false
        router.push(`/task/${taskDetail.value!.task_id}`)
      }, 2000)
    } else {
      alert(response.msg || '更新失败')
    }
  } catch (err: any) {
    alert(err.message || '网络错误，请重试')
  } finally {
    isSubmitting.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 组件挂载时获取任务详情
onMounted(() => {
  fetchTaskDetail()
})
</script>

<style scoped>
.edit-task-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background-color: #f8fafc;
  min-height: 100vh;
}

/* 加载和错误状态 */
.loading-state, .error-state, .permission-denied {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon, .permission-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.retry-btn, .back-btn {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background-color 0.2s;
}

.retry-btn:hover, .back-btn:hover {
  background-color: #2563eb;
}

/* 页面头部 */
.page-header {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-title h1 {
  margin: 0;
  color: #1f2937;
  font-size: 28px;
  font-weight: 700;
}

.header-subtitle {
  margin: 8px 0 0 0;
  color: #6b7280;
  font-size: 16px;
}

/* 表单样式 */
.edit-form-container {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 32px;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 24px;
}

.form-section:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
}

.section-icon {
  font-size: 20px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-weight: 500;
  color: #374151;
}

.form-input, .form-select, .form-textarea {
  padding: 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

/* 图片管理 */
.existing-images h4, .remark-images-section h4, .add-images-section h4 {
  margin: 16px 0 12px 0;
  color: #374151;
  font-weight: 500;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.image-item {
  position: relative;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.image-preview {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.image-actions {
  display: flex;
  gap: 8px;
  padding: 8px;
}

.action-btn {
  flex: 1;
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.update-btn {
  background-color: #f59e0b;
  color: white;
}

.update-btn:hover {
  background-color: #d97706;
}

.delete-btn {
  background-color: #ef4444;
  color: white;
}

.delete-btn:hover {
  background-color: #dc2626;
}

.image-order {
  padding: 4px 8px;
  background-color: #f3f4f6;
  font-size: 12px;
  color: #6b7280;
  text-align: center;
}

/* 图片更新表单 */
.update-image-form {
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.update-image-form h4 {
  margin: 0 0 16px 0;
  color: #374151;
}

/* 新图片和备注图片 */
.new-image-item, .remark-image-item {
  display: flex;
  gap: 16px;
  align-items: end;
  margin-bottom: 16px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.new-image-item .form-group, .remark-image-item .form-group {
  flex: 1;
}

.add-btn, .remove-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.add-btn {
  background-color: #10b981;
  color: white;
  margin-top: 12px;
}

.add-btn:hover {
  background-color: #059669;
}

.remove-btn {
  background-color: #ef4444;
  color: white;
  height: fit-content;
}

.remove-btn:hover {
  background-color: #dc2626;
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.submit-btn {
  background-color: #3b82f6;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  flex: 1;
  max-width: 200px;
}

.submit-btn:hover:not(:disabled) {
  background-color: #2563eb;
}

.submit-btn:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

.reset-btn {
  background-color: #6b7280;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.reset-btn:hover {
  background-color: #4b5563;
}

.confirm-btn, .cancel-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-right: 8px;
}

.confirm-btn {
  background-color: #10b981;
  color: white;
}

.confirm-btn:hover {
  background-color: #059669;
}

.cancel-btn {
  background-color: #6b7280;
  color: white;
}

.cancel-btn:hover {
  background-color: #4b5563;
}

/* 成功消息 */
.success-message {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  border: 1px solid #10b981;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  z-index: 1000;
}

.success-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .edit-task-container {
    padding: 16px;
  }
  
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .new-image-item, .remark-image-item {
    flex-direction: column;
    align-items: stretch;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .submit-btn {
    max-width: none;
  }
}
</style> 