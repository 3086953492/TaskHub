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

    <!-- 任务列表 -->
    <div class="task-list" v-if="paginatedTasks.length > 0">
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
      <p>没有找到符合条件的任务</p>
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
import { ref, computed, onMounted } from 'vue'
import TaskItem from './TaskItem.vue'

interface Task {
  task_id: number
  status: number
  created_at: string
  title: string
  priority: number
  due_date?: string
}

// 响应式数据
const tasks = ref<Task[]>([])
const currentPage = ref(1)
const pageSize = ref(6)
const statusFilter = ref('')
const priorityFilter = ref('')
const searchQuery = ref('')

// 模拟数据
const mockTasks: Task[] = [
  {
    task_id: 1,
    status: 1,
    created_at: '2024-01-15T09:30:00Z',
    title: '完成项目需求文档',
    priority: 1,
    due_date: '2024-01-25T18:00:00Z'
  },
  {
    task_id: 2,
    status: 2,
    created_at: '2024-01-14T14:20:00Z',
    title: '设计系统架构图',
    priority: 1,
    due_date: '2024-01-20T17:00:00Z'
  },
  {
    task_id: 3,
    status: 3,
    created_at: '2024-01-13T11:15:00Z',
    title: '编写API接口文档',
    priority: 2,
    due_date: '2024-01-18T16:00:00Z'
  },
  {
    task_id: 4,
    status: 1,
    created_at: '2024-01-12T16:45:00Z',
    title: '数据库表结构设计',
    priority: 2,
    due_date: '2024-01-22T15:30:00Z'
  },
  {
    task_id: 5,
    status: 2,
    created_at: '2024-01-11T10:00:00Z',
    title: '前端页面原型设计',
    priority: 1,
    due_date: '2024-01-19T12:00:00Z'
  },
  {
    task_id: 6,
    status: 4,
    created_at: '2024-01-10T13:30:00Z',
    title: '老版本系统维护',
    priority: 3,
    due_date: '2024-01-15T09:00:00Z'
  },
  {
    task_id: 7,
    status: 1,
    created_at: '2024-01-16T08:15:00Z',
    title: '用户体验优化方案',
    priority: 2,
    due_date: '2024-01-28T17:00:00Z'
  },
  {
    task_id: 8,
    status: 3,
    created_at: '2024-01-09T15:20:00Z',
    title: '性能测试报告',
    priority: 1,
    due_date: '2024-01-17T14:00:00Z'
  },
  {
    task_id: 9,
    status: 2,
    created_at: '2024-01-17T11:45:00Z',
    title: '移动端适配开发',
    priority: 2,
    due_date: '2024-02-01T18:00:00Z'
  },
  {
    task_id: 10,
    status: 1,
    created_at: '2024-01-18T09:00:00Z',
    title: '安全漏洞修复',
    priority: 1,
    due_date: '2024-01-21T16:00:00Z'
  }
]

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

// 生命周期
onMounted(() => {
  tasks.value = mockTasks
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