<template>
  <div class="task-list-container">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1>TaskHub 任务管理</h1>
      <p>高效管理您的任务和项目</p>
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
    <div class="pagination-info" v-if="tasks.length > 0">
      共 {{ tasks.length }} 条记录，每页 {{ pageSize }} 条，第 {{ currentPage }} / {{ totalPages }} 页
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
const totalPages = computed(() => {
  return Math.ceil(tasks.value.length / pageSize.value)
})

const paginatedTasks = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return tasks.value.slice(start, end)
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



// 方法

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
  min-height: calc(100vh - 64px); /* 减去导航栏高度 */
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
  
  .pagination {
    flex-direction: column;
    gap: 12px;
  }
  
  .page-numbers {
    order: -1;
  }
}
</style> 