<template>
  <div class="task-list-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>TaskHub 任务管理</h1>
      <p>高效管理您的任务和项目</p>
    </div>

    <!-- 筛选器 -->
    <div class="filters">
      <div class="filter-group">
        <label>状态筛选:</label>
        <select v-model="statusFilter" @change="applyFilters">
          <option value="">全部状态</option>
          <option value="1">待处理</option>
          <option value="2">进行中</option>
          <option value="3">已完成</option>
          <option value="4">已取消</option>
        </select>
      </div>

      <div class="filter-group">
        <label>优先级筛选:</label>
        <select v-model="priorityFilter" @change="applyFilters">
          <option value="">全部优先级</option>
          <option value="1">高优先级</option>
          <option value="2">中优先级</option>
          <option value="3">低优先级</option>
        </select>
      </div>

      <div class="filter-group">
        <label>搜索:</label>
        <input 
          type="text" 
          v-model="searchQuery" 
          @input="applyFilters"
          placeholder="搜索任务标题..."
          class="search-input"
        />
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="stats">
      <div class="stat-item">
        <span class="stat-number">{{ stats.total }}</span>
        <span class="stat-label">总任务数</span>
      </div>
      <div class="stat-item">
        <span class="stat-number">{{ stats.pending }}</span>
        <span class="stat-label">待处理</span>
      </div>
      <div class="stat-item">
        <span class="stat-number">{{ stats.inProgress }}</span>
        <span class="stat-label">进行中</span>
      </div>
      <div class="stat-item">
        <span class="stat-number">{{ stats.completed }}</span>
        <span class="stat-label">已完成</span>
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
    <div v-else-if="paginatedTasks.length > 0" class="task-list">
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
      <p>{{ isLoggedIn ? '没有找到符合条件的任务' : '请先登录查看任务' }}</p>
      <router-link v-if="!isLoggedIn" to="/login" class="login-link">前往登录</router-link>
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
          @click="currentPage = page"
          :class="['page-number', { active: page === currentPage }]"
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
import TaskItem from './TaskItem.vue'
import { getTaskList } from '../api/task'
import type { TaskListItem } from '../api/task'
import { useAuth } from '../composables/useAuth'

interface Task {
  task_id: number
  status: number
  created_at: string
  title: string
  priority: number
  due_date?: string
}

const { user, isLoggedIn } = useAuth()

// 响应式数据
const tasks = ref<Task[]>([])
const currentPage = ref(1)
const pageSize = ref(6)
const statusFilter = ref('')
const priorityFilter = ref('')
const searchQuery = ref('')
const isLoading = ref(false)
const error = ref('')

// 获取任务列表
const fetchTasks = async () => {
  if (!isLoggedIn.value) {
    error.value = '请先登录'
    return
  }

  isLoading.value = true
  error.value = ''

  try {
    // 根据用户角色构建请求参数
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    // 管理员获取所有任务，普通用户添加assignee_id=0条件
    if (user.value?.role !== 'admin') {
      params.assignee_id = 0
    }

    const response = await getTaskList(params)
    
    if (response.code === 200 && response.data) {
      tasks.value = response.data
    } else {
      error.value = response.msg || '获取任务列表失败'
    }
  } catch (err) {
    console.error('获取任务列表失败:', err)
    error.value = '网络错误，请稍后重试'
  } finally {
    isLoading.value = false
  }
}

// 计算属性
const filteredTasks = computed(() => {
  let filtered = [...tasks.value]
  
  // 状态筛选
  if (statusFilter.value) {
    filtered = filtered.filter(task => task.status.toString() === statusFilter.value)
  }
  
  // 优先级筛选
  if (priorityFilter.value) {
    filtered = filtered.filter(task => task.priority.toString() === priorityFilter.value)
  }
  
  // 搜索筛选
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(task => 
      task.title.toLowerCase().includes(query)
    )
  }
  
  return filtered
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

// 统计信息
const stats = computed(() => {
  const total = tasks.value.length
  const pending = tasks.value.filter(task => task.status === 1).length
  const inProgress = tasks.value.filter(task => task.status === 2).length
  const completed = tasks.value.filter(task => task.status === 3).length
  
  return { total, pending, inProgress, completed }
})

// 方法
const applyFilters = () => {
  currentPage.value = 1
}

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
watch([currentPage, isLoggedIn], () => {
  if (isLoggedIn.value) {
    fetchTasks()
  }
}, { immediate: true })

// 监听用户变化，重新获取数据
watch(user, () => {
  if (user.value) {
    currentPage.value = 1
    fetchTasks()
  }
})

// 生命周期
onMounted(() => {
  if (isLoggedIn.value) {
    fetchTasks()
  }
})
</script>

<style scoped>
.task-list-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  background: #f8fafc;
  min-height: 100vh;
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

.filters {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  flex-wrap: wrap;
  justify-content: flex-start;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 200px;
  max-width: 250px;
  flex: 0 1 auto;
}

.filter-group label {
  font-weight: 500;
  color: #374151;
  font-size: 14px;
}

.filter-group select,
.search-input {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.filter-group select:focus,
.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 280px));
  gap: 16px;
  margin-bottom: 24px;
  justify-content: center;
}

.stat-item {
  background: white;
  padding: 20px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.stat-number {
  display: block;
  font-size: 32px;
  font-weight: 700;
  color: #3b82f6;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.task-list {
  margin-bottom: 32px;
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

.login-link {
  display: inline-block;
  background: #3b82f6;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  padding: 10px 20px;
  margin-top: 16px;
  transition: background-color 0.2s;
}

.login-link:hover {
  background: #2563eb;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 64px;
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
  color: white;
  border-color: #3b82f6;
}

.pagination-info {
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

@media (max-width: 768px) {
  .task-list-container {
    padding: 16px;
  }
  
  .filters {
    flex-direction: column;
    gap: 16px;
  }
  
  .filter-group {
    min-width: auto;
  }
  
  .stats {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .pagination {
    flex-direction: column;
    gap: 12px;
  }
  
  .page-numbers {
    order: -1;
  }
}
</style> 