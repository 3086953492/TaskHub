<template>
  <div class="task-history-container">
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载历史记录...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">❌</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button @click="fetchHistory" class="retry-btn">重试</button>
    </div>

    <!-- 历史记录内容 -->
    <div v-else class="history-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-actions">
          <button @click="goBack" class="back-btn">
            <span class="btn-icon">←</span>
            <span>返回详情</span>
          </button>
          
          <div class="header-info">
            <h1>任务历史记录</h1>
            <p class="task-info">任务 #{{ taskId }}</p>
          </div>
          
          <div class="header-stats">
            <div class="stat-item">
              <span class="stat-value">{{ historyList.length }}</span>
              <span class="stat-label">操作记录</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 历史记录列表 -->
      <div class="history-timeline">
        <!-- 空状态 -->
        <div v-if="historyList.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <h3>暂无历史记录</h3>
          <p>该任务还没有任何操作记录</p>
        </div>

        <!-- 历史记录项 -->
        <div v-else class="timeline-list">
          <div 
            v-for="(item, index) in historyList" 
            :key="item.id"
            class="timeline-item"
            :class="{ 'is-last': index === historyList.length - 1 }"
          >
            <!-- 时间线连接器 -->
            <div class="timeline-connector">
              <div class="timeline-dot" :style="{ backgroundColor: getActionColor(item.action) }">
                <span class="action-icon">{{ getActionIcon(item.action) }}</span>
              </div>
              <div v-if="index < historyList.length - 1" class="timeline-line"></div>
            </div>

            <!-- 历史记录内容 -->
            <div class="timeline-content">
              <div class="history-card">
                <div class="card-header">
                  <div class="action-info">
                    <h3 class="action-title">{{ item.action }}</h3>
                    <div class="action-meta">
                      <span class="operator">操作人: #{{ item.operator_id }}</span>
                      <span class="timestamp">{{ formatDate(item.created_at) }}</span>
                    </div>
                  </div>
                  
                  <div class="action-badge" :style="{ backgroundColor: getActionColor(item.action) }">
                    {{ getActionType(item.action) }}
                  </div>
                </div>

                <!-- 备注信息 -->
                <div v-if="item.remark" class="card-content">
                  <div class="remark-section">
                    <div class="remark-label">
                      <span class="remark-icon">💬</span>
                      <span>备注信息</span>
                    </div>
                    <p class="remark-text">{{ item.remark }}</p>
                  </div>
                </div>

                <!-- 相关图片 -->
                <div v-if="item.images && item.images.length > 0" class="card-images">
                  <div class="images-label">
                    <span class="images-icon">🖼️</span>
                    <span>相关图片 ({{ item.images.length }})</span>
                  </div>
                  <div class="images-grid">
                    <div 
                      v-for="(imageUrl, imgIndex) in item.images" 
                      :key="imgIndex"
                      class="image-item"
                      @click="previewImage(imageUrl)"
                    >
                      <img :src="imageUrl" :alt="`历史图片 ${imgIndex + 1}`" class="history-image" />
                      <div class="image-overlay">
                        <span class="preview-icon">🔍</span>
                      </div>
                    </div>
                  </div>
                </div>
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
import { getTaskHistory } from '../api/task'
import type { TaskHistory } from '../api/task'

const route = useRoute()
const router = useRouter()

// 状态
const isLoading = ref(false)
const error = ref('')
const historyList = ref<TaskHistory[]>([])
const previewImageUrl = ref('')

// 获取任务ID
const taskId = Number(route.params.id)

// 获取历史记录
const fetchHistory = async () => {
  if (!taskId || isNaN(taskId)) {
    error.value = '无效的任务ID'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    const response = await getTaskHistory(taskId)
    
    if (response.code === 200) {
      // 确保数据是数组，如果为null或undefined则设为空数组
      historyList.value = Array.isArray(response.data) ? response.data : []
    } else {
      error.value = response.msg || '获取历史记录失败'
    }
  } catch (err) {
    console.error('获取历史记录失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return `今天 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else if (diffDays === 1) {
    return `昨天 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else if (diffDays < 7) {
    return `${diffDays}天前 ${date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })}`
  } else {
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }
}

// 获取操作图标
const getActionIcon = (action: string): string => {
  const iconMap: Record<string, string> = {
    '创建任务': '✨',
    '分配任务': '👤',
    '更新任务': '📝',
    '完成任务': '✅',
    '取消任务': '❌',
    '开始任务': '🚀',
    '暂停任务': '⏸️',
    '删除任务': '🗑️',
    '修改状态': '🔄'
  }
  return iconMap[action] || '📋'
}

// 获取操作颜色
const getActionColor = (action: string): string => {
  const colorMap: Record<string, string> = {
    '创建任务': '#10b981',
    '分配任务': '#3b82f6',
    '更新任务': '#f59e0b',
    '完成任务': '#10b981',
    '取消任务': '#ef4444',
    '开始任务': '#8b5cf6',
    '暂停任务': '#6b7280',
    '删除任务': '#ef4444',
    '修改状态': '#06b6d4'
  }
  return colorMap[action] || '#6b7280'
}

// 获取操作类型
const getActionType = (action: string): string => {
  if (action.includes('创建')) return '创建'
  if (action.includes('分配')) return '分配'
  if (action.includes('更新') || action.includes('修改')) return '更新'
  if (action.includes('完成')) return '完成'
  if (action.includes('取消')) return '取消'
  if (action.includes('开始')) return '开始'
  if (action.includes('暂停')) return '暂停'
  if (action.includes('删除')) return '删除'
  return '操作'
}

// 预览图片
const previewImage = (imageUrl: string) => {
  previewImageUrl.value = imageUrl
}

// 关闭图片预览
const closeImagePreview = () => {
  previewImageUrl.value = ''
}

// 返回任务详情页
const goBack = () => {
  // 使用 replace 替换当前历史记录，避免循环导航
  router.replace(`/task/${taskId}`)
}

// 生命周期
onMounted(() => {
  fetchHistory()
})
</script>

<style scoped>
.task-history-container {
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
  background: #f8fafc;
  min-height: calc(100vh - 64px);
}

.loading-state,
.error-state {
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

.loading-state p,
.error-state p {
  color: #6b7280;
  margin: 0;
}

.error-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-state h3 {
  color: #dc2626;
  margin: 0 0 8px 0;
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
  margin-top: 20px;
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
  gap: 20px;
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  background: #f3f4f6;
  color: #374151;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.back-btn:hover {
  background: #e5e7eb;
}

.btn-icon {
  font-size: 14px;
}

.header-info {
  text-align: center;
  flex: 1;
}

.header-info h1 {
  color: #1f2937;
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 4px 0;
}

.task-info {
  color: #6b7280;
  font-size: 14px;
  margin: 0;
}

.header-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  text-align: center;
  background: #f8fafc;
  padding: 12px 16px;
  border-radius: 8px;
}

.stat-value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

.history-timeline {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.empty-state h3 {
  color: #1f2937;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #6b7280;
  margin: 0;
}

.timeline-list {
  position: relative;
}

.timeline-item {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
}

.timeline-item.is-last {
  margin-bottom: 0;
}

.timeline-connector {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
}

.timeline-dot {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1;
}

.timeline-line {
  width: 2px;
  flex: 1;
  background: #e5e7eb;
  margin-top: 8px;
  min-height: 20px;
}

.timeline-content {
  flex: 1;
  min-width: 0;
}

.history-card {
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
  transition: all 0.2s;
}

.history-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  background: white;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.action-info {
  flex: 1;
}

.action-title {
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 8px 0;
}

.action-meta {
  display: flex;
  gap: 16px;
  font-size: 14px;
  color: #6b7280;
}

.action-badge {
  color: white;
  font-size: 12px;
  font-weight: 600;
  padding: 4px 12px;
  border-radius: 16px;
  white-space: nowrap;
}

.card-content {
  padding: 0 20px 20px;
}

.remark-section {
  background: white;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #e5e7eb;
}

.remark-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-weight: 500;
  font-size: 14px;
  margin-bottom: 8px;
}

.remark-icon {
  font-size: 16px;
}

.remark-text {
  color: #1f2937;
  line-height: 1.6;
  margin: 0;
}

.card-images {
  padding: 0 20px 20px;
}

.images-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-weight: 500;
  font-size: 14px;
  margin-bottom: 12px;
}

.images-icon {
  font-size: 16px;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.image-item {
  position: relative;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s;
  border: 2px solid #e5e7eb;
}

.image-item:hover {
  transform: scale(1.05);
}

.history-image {
  width: 100%;
  height: 80px;
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
  font-size: 16px;
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
  .task-history-container {
    padding: 16px;
  }
  
  .page-header {
    padding: 20px;
  }
  
  .header-actions {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .header-info {
    order: -1;
  }
  
  .header-info h1 {
    font-size: 20px;
  }
  
  .timeline-item {
    gap: 12px;
  }
  
  .timeline-dot {
    width: 32px;
    height: 32px;
    font-size: 14px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .action-meta {
    flex-direction: column;
    gap: 4px;
  }
  
  .images-grid {
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  }
  
  .history-image {
    height: 60px;
  }
}
</style> 