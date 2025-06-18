<template>
  <div class="task-item">
    <div class="task-header">
      <h3 class="task-title clickable" @click="goToDetail">{{ task.title }}</h3>
      <div class="task-badges">
        <span class="status-badge" :class="statusClass">
          {{ statusText }}
        </span>
        <span class="priority-badge" :class="priorityClass">
          {{ priorityText }}
        </span>
      </div>
    </div>
    
    <div class="task-meta">
      <div class="meta-item">
        <span class="label">创建时间:</span>
        <span class="value">{{ formattedCreatedAt }}</span>
      </div>
      <div class="meta-item" v-if="task.due_date">
        <span class="label">截止时间:</span>
        <span class="value" :class="{ 'overdue': isOverdue }">{{ formattedDueDate }}</span>
      </div>
      <div class="meta-item">
        <span class="label">任务ID:</span>
        <span class="value">#{{ task.task_id }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

interface Task {
  task_id: number
  status: number
  created_at: string
  title: string
  priority: number
  due_date?: string
}

const props = defineProps<{
  task: Task
}>()

const router = useRouter()

// 跳转到任务详情页
const goToDetail = () => {
  router.push(`/task/${props.task.task_id}`)
}

// 状态文本映射
const statusText = computed(() => {
  const statusMap = {
    1: '待处理',
    2: '进行中',
    3: '已完成',
    4: '已取消'
  }
  return statusMap[props.task.status as keyof typeof statusMap] || '未知'
})

// 状态样式类
const statusClass = computed(() => {
  const classMap = {
    1: 'status-pending',
    2: 'status-progress',
    3: 'status-completed',
    4: 'status-cancelled'
  }
  return classMap[props.task.status as keyof typeof classMap] || ''
})

// 优先级文本映射
const priorityText = computed(() => {
  const priorityMap = {
    1: '高优先级',
    2: '中优先级',
    3: '低优先级'
  }
  return priorityMap[props.task.priority as keyof typeof priorityMap] || '未知'
})

// 优先级样式类
const priorityClass = computed(() => {
  const classMap = {
    1: 'priority-high',
    2: 'priority-medium',
    3: 'priority-low'
  }
  return classMap[props.task.priority as keyof typeof classMap] || ''
})

// 格式化时间
const formattedCreatedAt = computed(() => {
  return new Date(props.task.created_at).toLocaleString('zh-CN')
})

const formattedDueDate = computed(() => {
  if (!props.task.due_date) return ''
  return new Date(props.task.due_date).toLocaleString('zh-CN')
})

// 检查是否已过期
const isOverdue = computed(() => {
  if (!props.task.due_date || props.task.status === 3) return false
  return new Date(props.task.due_date) < new Date()
})
</script>

<style scoped>
.task-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
}

.task-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.task-title {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
  flex: 1;
  margin-right: 16px;
}

.task-title.clickable {
  cursor: pointer;
  transition: color 0.2s ease;
}

.task-title.clickable:hover {
  color: #3b82f6;
}

.task-badges {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.status-badge, .priority-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-pending {
  background: #fef3c7;
  color: #92400e;
}

.status-progress {
  background: #dbeafe;
  color: #1e40af;
}

.status-completed {
  background: #dcfce7;
  color: #166534;
}

.status-cancelled {
  background: #fee2e2;
  color: #dc2626;
}

.priority-high {
  background: #fee2e2;
  color: #dc2626;
}

.priority-medium {
  background: #fef3c7;
  color: #d97706;
}

.priority-low {
  background: #f3f4f6;
  color: #6b7280;
}

.task-meta {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
  font-size: 14px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label {
  color: #6b7280;
  font-weight: 500;
}

.value {
  color: #1f2937;
  font-weight: 400;
}

.value.overdue {
  color: #dc2626;
  font-weight: 600;
}

@media (max-width: 768px) {
  .task-header {
    flex-direction: column;
    gap: 12px;
  }
  
  .task-title {
    margin-right: 0;
  }
  
  .task-meta {
    flex-direction: column;
    gap: 8px;
  }
}
</style> 