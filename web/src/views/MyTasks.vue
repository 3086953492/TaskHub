<template>
  <div class="my-tasks-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <div class="title-section">
          <h1>我的任务</h1>
          <p>查看和管理您的个人任务</p>
        </div>
        <div class="action-section">
          <router-link to="/task/create" class="create-task-btn">
            <span class="btn-icon">➕</span>
            <span>创建任务</span>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 选项卡 -->
    <div class="tab-container">
      <div class="tab-buttons">
        <button 
          @click="activeTab = 'assigned'" 
          :class="['tab-button', { active: activeTab === 'assigned' }]"
        >
          <span class="tab-icon">📋</span>
          <span>我认领的任务</span>
          <span v-if="assignedCount > 0" class="task-count">{{ assignedCount }}</span>
        </button>
        <button 
          @click="activeTab = 'created'" 
          :class="['tab-button', { active: activeTab === 'created' }]"
        >
          <span class="tab-icon">✏️</span>
          <span>我发布的任务</span>
          <span v-if="createdCount > 0" class="task-count">{{ createdCount }}</span>
        </button>
      </div>
    </div>

    <!-- 状态筛选 -->
    <div class="filter-container">
      <div class="filter-section">
        <label for="status-filter">状态筛选：</label>
        <select id="status-filter" v-model="statusFilter" class="status-select">
          <option :value="0">全部状态</option>
          <option v-for="status in STATUS_OPTIONS" :key="status.value" :value="status.value">
            {{ status.label }}
          </option>
        </select>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-state">
      <div class="loading-spinner"></div>
      <p>正在加载任务列表...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="error-state">
      <div class="error-icon">❌</div>
      <h3>加载失败</h3>
      <p>{{ error }}</p>
      <button @click="fetchTasks" class="retry-btn">重试</button>
    </div>

    <!-- 任务列表 -->
    <div v-else-if="filteredTasks.length > 0" class="task-list">
      <TaskItem 
        v-for="task in paginatedTasks" 
        :key="task.task_id" 
        :task="task"
      />
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <div class="empty-icon">📋</div>
      <h3>暂无任务</h3>
      <p v-if="activeTab === 'assigned'">您还没有认领任何任务</p>
      <p v-else>您还没有发布任何任务</p>
      <router-link v-if="activeTab === 'created'" to="/task/create" class="create-link">
        创建第一个任务
      </router-link>
    </div>

    <!-- 分页器 -->
    <div class="pagination" v-if="totalPages > 1">
      <button 
        @click="prevPage" 
        :disabled="currentPage === 1"
        class="page-btn"
      >
        上一页
      </button>
      
      <div class="page-numbers">
        <button 
          v-for="page in visiblePages" 
          :key="page"
          @click="typeof page === 'number' && (currentPage = page)"
          :class="['page-number', { active: page === currentPage }]"
          :disabled="typeof page === 'string'"
        >
          {{ page }}
        </button>
      </div>
      
      <button 
        @click="nextPage" 
        :disabled="currentPage === totalPages"
        class="page-btn"
      >
        下一页
      </button>
    </div>

    <!-- 分页信息 -->
    <div class="pagination-info" v-if="filteredTasks.length > 0">
      共 {{ filteredTasks.length }} 条记录，每页 {{ pageSize }} 条，第 {{ currentPage }} / {{ totalPages }} 页
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import TaskItem from '../components/TaskItem.vue'
import { getTaskList, STATUS_OPTIONS, type TaskListItem } from '../api/task'
import { useAuth } from '../composables/useAuth'

const { user } = useAuth()

// 响应式数据
const activeTab = ref<'assigned' | 'created'>('assigned')
const assignedTasks = ref<TaskListItem[]>([])
const createdTasks = ref<TaskListItem[]>([])
const currentPage = ref(1)
const pageSize = ref(6)
const isLoading = ref(false)
const error = ref('')
const statusFilter = ref(0)

// 计算属性
const assignedCount = computed(() => assignedTasks.value.length)
const createdCount = computed(() => createdTasks.value.length)

const currentTasks = computed(() => {
  return activeTab.value === 'assigned' ? assignedTasks.value : createdTasks.value
})

const filteredTasks = computed(() => {
  if (statusFilter.value === 0) {
    return currentTasks.value
  }
  return currentTasks.value.filter(task => task.status === statusFilter.value)
})

const totalPages = computed(() => {
  return Math.ceil(filteredTasks.value.length / pageSize.value)
})

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTasks.value.slice(start, end)
})

const visiblePages = computed(() => {
  const delta = 2
  const range = []
  const rangeWithDots = []
  
  for (let i = Math.max(2, currentPage.value - delta); 
       i <= Math.min(totalPages.value - 1, currentPage.value + delta); 
       i++) {
    range.push(i)
  }
  
  if (currentPage.value - delta > 2) {
    rangeWithDots.push(1, '...')
  } else {
    rangeWithDots.push(1)
  }
  
  rangeWithDots.push(...range)
  
  if (currentPage.value + delta < totalPages.value - 1) {
    rangeWithDots.push('...', totalPages.value)
  } else {
    rangeWithDots.push(totalPages.value)
  }
  
  return [...new Set(rangeWithDots)]
})

// 获取任务列表
const fetchTasks = async () => {
  if (!user.value) {
    error.value = '请先登录'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    // 并行获取认领的任务和发布的任务
    const [assignedResponse, createdResponse] = await Promise.all([
      // 获取我认领的任务
      getTaskList({
        page: 1,
        page_size: 100, // 获取较多数据，前端分页
        assignee_id: user.value.id
      }),
      // 获取我发布的任务
      getTaskList({
        page: 1,
        page_size: 100,
        creator_id: user.value.id
      })
    ])
    
    if (assignedResponse.code === 200 && assignedResponse.data) {
      assignedTasks.value = assignedResponse.data
    } else {
      console.error('获取认领任务失败:', assignedResponse.msg)
    }

    if (createdResponse.code === 200 && createdResponse.data) {
      createdTasks.value = createdResponse.data
    } else {
      console.error('获取发布任务失败:', createdResponse.msg)
    }

  } catch (err) {
    console.error('获取任务列表失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 分页方法
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// 监听器
watch(activeTab, () => {
  currentPage.value = 1
})

watch(statusFilter, () => {
  currentPage.value = 1
})

// 生命周期
onMounted(() => {
  fetchTasks()
})
</script>

<style scoped>
.my-tasks-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: #f8fafc;
  min-height: calc(100vh - 64px);
}

.page-header {
  margin-bottom: 32px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}

.title-section {
  text-align: center;
  flex: 1;
}

.title-section h1 {
  color: #1f2937;
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.title-section p {
  color: #6b7280;
  font-size: 16px;
  margin: 0;
}

.action-section {
  flex-shrink: 0;
}

.create-task-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: #10b981;
  color: white;
  text-decoration: none;
  border-radius: 10px;
  font-weight: 500;
  transition: all 0.2s;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
}

.create-task-btn:hover {
  background: #059669;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.btn-icon {
  font-size: 14px;
}

/* 选项卡样式 */
.tab-container {
  margin-bottom: 24px;
}

.tab-buttons {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  gap: 4px;
}

.tab-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 20px;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  transition: all 0.2s;
  position: relative;
}

.tab-button:hover {
  background: #f3f4f6;
  color: #374151;
}

.tab-button.active {
  background: #3b82f6;
  color: white;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
}

.tab-icon {
  font-size: 16px;
}

.task-count {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 12px;
  font-weight: 600;
  min-width: 20px;
  text-align: center;
}

.tab-button.active .task-count {
  background: rgba(255, 255, 255, 0.2);
}

/* 筛选器样式 */
.filter-container {
  margin-bottom: 24px;
}

.filter-section {
  display: flex;
  align-items: center;
  gap: 12px;
  background: white;
  padding: 16px 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filter-section label {
  font-weight: 500;
  color: #374151;
  white-space: nowrap;
}

.status-select {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 14px;
  cursor: pointer;
  transition: border-color 0.2s;
}

.status-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.task-list {
  margin-bottom: 32px;
}

/* 加载、错误、空状态样式 */
.loading-state, .error-state, .empty-state {
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

.error-icon, .empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
}

.error-state h3, .empty-state h3 {
  color: #1f2937;
  margin: 0 0 8px 0;
}

.error-state p, .empty-state p, .loading-state p {
  color: #6b7280;
  margin: 0 0 20px 0;
}

.retry-btn, .create-link {
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  cursor: pointer;
  font-size: 14px;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.retry-btn:hover, .create-link:hover {
  background: #2563eb;
}

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.page-btn {
  padding: 8px 16px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-numbers {
  display: flex;
  gap: 4px;
}

.page-number {
  padding: 8px 12px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
  min-width: 40px;
}

.page-number:hover {
  background: #f3f4f6;
  border-color: #9ca3af;
}

.page-number.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.page-number:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-info {
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .my-tasks-container {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }

  .title-section {
    text-align: center;
  }

  .tab-buttons {
    flex-direction: column;
    gap: 8px;
  }

  .tab-button {
    justify-content: center;
  }

  .filter-section {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .filter-section label {
    text-align: center;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 8px;
  }

  .page-numbers {
    flex-wrap: wrap;
  }
}
</style> 