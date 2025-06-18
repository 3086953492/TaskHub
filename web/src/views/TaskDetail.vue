<template>
  <div class="task-detail-container">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载任务详情...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">❌</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button @click="fetchTaskDetail" class="retry-btn">重试</button>
    </div>

    <!-- 任务详情内容 -->
    <div v-else-if="taskDetail" class="task-detail-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-actions">
          <button @click="goBack" class="back-btn">
            <span class="btn-icon">←</span>
            <span>返回</span>
          </button>
          
          <div class="header-buttons">
            <button @click="viewHistory" class="history-btn">
              <span class="btn-icon">📋</span>
              <span>查看历史</span>
            </button>
          </div>
        </div>
        
        <div class="task-title-section">
          <h1>{{ taskDetail.title }}</h1>
          <div class="task-meta">
            <span class="task-id">任务 #{{ taskDetail.task_id }}</span>
            <span 
              class="status-badge" 
              :style="{ backgroundColor: getStatusOption(taskDetail.status).color }"
            >
              {{ getStatusOption(taskDetail.status).label }}
            </span>
            <span 
              class="priority-badge" 
              :style="{ backgroundColor: getPriorityOption(taskDetail.priority).color }"
            >
              {{ getPriorityOption(taskDetail.priority).label }}
            </span>
          </div>
        </div>
      </div>

      <!-- 任务详情卡片 -->
      <div class="detail-cards">
        <!-- 基本信息卡片 -->
        <div class="detail-card">
          <h3>基本信息</h3>
          <div class="info-grid">
            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">📝</span>
                <span>任务描述</span>
              </div>
              <div class="info-value">
                <p v-if="taskDetail.description" class="description-text">
                  {{ taskDetail.description }}
                </p>
                <span v-else class="empty-value">暂无描述</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">📅</span>
                <span>截止日期</span>
              </div>
              <div class="info-value">
                <span v-if="taskDetail.due_date" :class="getDueDateClass()">
                  {{ formatDate(taskDetail.due_date) }}
                </span>
                <span v-else class="empty-value">无截止日期</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">⏰</span>
                <span>创建时间</span>
              </div>
              <div class="info-value">
                {{ formatDate(taskDetail.created_at) }}
              </div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">🔄</span>
                <span>更新时间</span>
              </div>
              <div class="info-value">
                {{ formatDate(taskDetail.updated_at) }}
              </div>
            </div>

            <div v-if="taskDetail.completed_at" class="info-item">
              <div class="info-label">
                <span class="info-icon">✅</span>
                <span>完成时间</span>
              </div>
              <div class="info-value">
                {{ formatDate(taskDetail.completed_at) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 任务图片卡片 -->
        <div v-if="taskDetail.images && taskDetail.images.length > 0" class="detail-card">
          <h3>任务图片</h3>
          <div class="images-grid">
            <div 
              v-for="(image, index) in taskDetail.images" 
              :key="index"
              class="image-item"
              @click="previewImage(image.url)"
            >
              <img :src="image.url" :alt="`任务图片 ${index + 1}`" class="task-image" />
              <div class="image-overlay">
                <span class="preview-icon">🔍</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 任务统计卡片 -->
        <div class="detail-card stats-card">
          <h3>任务统计</h3>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-icon">👤</div>
              <div class="stat-content">
                <div class="stat-label">创建人ID</div>
                <div class="stat-value">#{{ taskDetail.creator_id }}</div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">🎯</div>
              <div class="stat-content">
                <div class="stat-label">分配人ID</div>
                <div class="stat-value">
                  {{ taskDetail.assignee_id ? `#${taskDetail.assignee_id}` : '未分配' }}
                </div>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">⏱️</div>
              <div class="stat-content">
                <div class="stat-label">运行时长</div>
                <div class="stat-value">{{ getTaskDuration() }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图片预览模态框 -->
    <div v-if="previewImageUrl" class="image-modal" @click="closeImagePreview">
      <div class="modal-content" @click.stop>
        <img :src="previewImageUrl" alt="图片预览" class="preview-image" />
        <button @click="closeImagePreview" class="close-modal-btn">×</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getTaskDetail, getPriorityOption, getStatusOption } from '../api/task'
import type { TaskDetail } from '../api/task'
import { message } from '../utils/message'

const route = useRoute()
const router = useRouter()

// 状态
const isLoading = ref(false)
const error = ref('')
const taskDetail = ref<TaskDetail | null>(null)
const previewImageUrl = ref('')

// 获取任务ID
const taskId = Number(route.params.id)

// 获取任务详情
const fetchTaskDetail = async () => {
  if (!taskId || isNaN(taskId)) {
    error.value = '无效的任务ID'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    const response = await getTaskDetail(taskId)
    
    if (response.code === 200 && response.data) {
      taskDetail.value = response.data
    } else {
      error.value = response.msg || '获取任务详情失败'
    }
  } catch (err) {
    console.error('获取任务详情失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取截止日期样式
const getDueDateClass = (): string => {
  if (!taskDetail.value?.due_date) return ''
  
  const dueDate = new Date(taskDetail.value.due_date)
  const now = new Date()
  const diffDays = Math.ceil((dueDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  
  if (diffDays < 0) return 'overdue'
  if (diffDays <= 1) return 'urgent'
  if (diffDays <= 3) return 'warning'
  return 'normal'
}

// 计算任务运行时长
const getTaskDuration = (): string => {
  if (!taskDetail.value) return '未知'
  
  const startDate = new Date(taskDetail.value.created_at)
  const endDate = taskDetail.value.completed_at 
    ? new Date(taskDetail.value.completed_at)
    : new Date()
  
  const diffMs = endDate.getTime() - startDate.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  
  if (diffDays > 0) {
    return `${diffDays}天${diffHours}小时`
  } else if (diffHours > 0) {
    return `${diffHours}小时`
  } else {
    return '不到1小时'
  }
}

// 预览图片
const previewImage = (imageUrl: string) => {
  previewImageUrl.value = imageUrl
}

// 关闭图片预览
const closeImagePreview = () => {
  previewImageUrl.value = ''
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 查看历史记录
const viewHistory = () => {
  message.info('功能开发中', '历史记录功能正在开发中，敬请期待')
}

// 生命周期
onMounted(() => {
  fetchTaskDetail()
})
</script>

<style scoped>
.task-detail-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: #f8fafc;
  min-height: calc(100vh - 64px);
}

.loading-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-state p {
  color: #6b7280;
  margin: 0;
}

.error-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-state h3 {
  color: #dc2626;
  margin: 0 0 8px 0;
}

.error-state p {
  color: #6b7280;
  margin: 0 0 20px 0;
}

.retry-btn {
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.retry-btn:hover {
  background: #2563eb;
}

.page-header {
  background: white;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.back-btn,
.history-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn {
  background: #f3f4f6;
  color: #374151;
}

.back-btn:hover {
  background: #e5e7eb;
}

.history-btn {
  background: #3b82f6;
  color: white;
}

.history-btn:hover {
  background: #2563eb;
}

.btn-icon {
  font-size: 14px;
}

.task-title-section h1 {
  color: #1f2937;
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 12px 0;
}

.task-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.task-id {
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
}

.status-badge,
.priority-badge {
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 16px;
}

.detail-cards {
  display: grid;
  gap: 24px;
}

.detail-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.detail-card h3 {
  color: #1f2937;
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 20px 0;
}

.info-grid {
  display: grid;
  gap: 20px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-weight: 500;
  min-width: 120px;
  flex-shrink: 0;
}

.info-icon {
  font-size: 16px;
}

.info-value {
  flex: 1;
  color: #1f2937;
}

.description-text {
  margin: 0;
  line-height: 1.6;
}

.empty-value {
  color: #9ca3af;
  font-style: italic;
}

.normal { color: #10b981; }
.warning { color: #f59e0b; }
.urgent { color: #ef4444; }
.overdue { color: #dc2626; font-weight: 600; }

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.image-item {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s;
  border: 2px solid #e5e7eb;
}

.image-item:hover {
  transform: scale(1.02);
}

.task-image {
  width: 100%;
  height: 150px;
  object-fit: cover;
  display: block;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.image-item:hover .image-overlay {
  opacity: 1;
}

.preview-icon {
  color: white;
  font-size: 24px;
}

.stats-card .stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 12px;
}

.stat-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 8px;
}

.stat-content {
  flex: 1;
}

.stat-label {
  color: #6b7280;
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 4px;
}

.stat-value {
  color: #1f2937;
  font-weight: 600;
}

.image-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
}

.preview-image {
  max-width: 100%;
  max-height: 100%;
  border-radius: 8px;
}

.close-modal-btn {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 32px;
  height: 32px;
  background: white;
  border: none;
  border-radius: 50%;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

@media (max-width: 768px) {
  .task-detail-container {
    padding: 16px;
  }
  
  .page-header {
    padding: 20px;
  }
  
  .header-actions {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .header-buttons {
    display: flex;
    justify-content: center;
  }
  
  .task-title-section h1 {
    font-size: 24px;
  }
  
  .task-meta {
    justify-content: center;
  }
  
  .info-item {
    flex-direction: column;
    gap: 8px;
  }
  
  .info-label {
    min-width: auto;
  }
  
  .images-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  }
  
  .task-image {
    height: 120px;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style> 