<template>
  <div class="user-profile-container">
    <div class="page-header">
      <h1>用户信息</h1>
      <p>管理您的个人信息和账户设置</p>
    </div>

    <!-- 用户信息卡片 -->
    <div class="profile-content">
      <div class="profile-card">
        <!-- 头像区域 -->
        <div class="avatar-section">
          <div class="avatar-container">
            <img 
              v-if="user?.avatar" 
              :src="getImageUrl(user.avatar)" 
              :alt="user.nickname || user.username"
              class="avatar-image"
            />
            <div v-else class="avatar-placeholder">
              {{ getAvatarText(user?.nickname || user?.username || '') }}
            </div>
          </div>
          <div class="user-title">
            <h2>{{ user?.nickname || user?.username }}</h2>
            <span class="user-role" :class="roleClass">{{ roleText }}</span>
          </div>
        </div>

        <!-- 用户信息详情 -->
        <div class="info-section">
          <div class="info-grid">
            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">👤</span>
                <span>用户名</span>
              </div>
              <div class="info-value">{{ user?.username || '-' }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">✏️</span>
                <span>昵称</span>
              </div>
              <div class="info-value">{{ user?.nickname || '-' }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">📧</span>
                <span>邮箱</span>
              </div>
              <div class="info-value">{{ user?.email || '-' }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">🔑</span>
                <span>用户角色</span>
              </div>
              <div class="info-value">
                <span class="role-badge" :class="roleClass">{{ roleText }}</span>
              </div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">🆔</span>
                <span>用户ID</span>
              </div>
              <div class="info-value">#{{ user?.id || '-' }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">
                <span class="info-icon">⏰</span>
                <span>登录状态</span>
              </div>
              <div class="info-value">
                <span class="status-badge online">在线</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-section">
          <button class="action-btn primary" @click="editProfile">
            <span class="btn-icon">✏️</span>
            <span>编辑资料</span>
          </button>
          
          <button class="action-btn danger" @click="handleLogout">
            <span class="btn-icon">🚪</span>
            <span>退出登录</span>
          </button>
        </div>
      </div>

      <!-- 统计信息卡片 -->
      <div class="stats-card">
        <h3>账户统计</h3>
        <div class="stats-grid">
          <div class="stat-item">
            <div class="stat-number">{{ user?.role === 'admin' ? '∞' : '?' }}</div>
            <div class="stat-label">可管理任务</div>
          </div>
          
          <div class="stat-item">
            <div class="stat-number">{{ getJoinDays() }}</div>
            <div class="stat-label">加入天数</div>
          </div>
          
          <div class="stat-item">
            <div class="stat-number">1</div>
            <div class="stat-label">当前会话</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { getImageUrl } from '../api/upload'

const router = useRouter()
const { user, logout } = useAuth()

// 用户角色显示
const roleText = computed(() => {
  return user.value?.role === 'admin' ? '管理员' : '普通用户'
})

const roleClass = computed(() => {
  return user.value?.role === 'admin' ? 'role-admin' : 'role-user'
})

// 获取头像文字
const getAvatarText = (name: string) => {
  return name.charAt(0).toUpperCase()
}

// 计算加入天数（模拟）
const getJoinDays = () => {
  // 这里可以根据用户注册时间计算，暂时返回模拟数据
  return Math.floor(Math.random() * 365) + 1
}

// 事件处理
const editProfile = () => {
  router.push('/profile/edit')
}

const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    logout()
    router.push('/login')
  }
}
</script>

<style scoped>
.user-profile-container {
  width: 100%;
  max-width: 1400px;
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

.profile-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.profile-card,
.stats-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.profile-card {
  padding: 32px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 32px;
  padding-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.avatar-container {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
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

.user-title h2 {
  margin: 0 0 8px 0;
  color: #1f2937;
  font-size: 24px;
  font-weight: 700;
}

.user-role {
  font-size: 14px;
  padding: 4px 12px;
  border-radius: 20px;
  font-weight: 500;
}

.role-admin {
  background: #fee2e2;
  color: #dc2626;
}

.role-user {
  background: #e0e7ff;
  color: #3730a3;
}

.info-section {
  margin-bottom: 32px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0;
  border-bottom: 1px solid #f3f4f6;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-weight: 500;
}

.info-icon {
  font-size: 16px;
}

.info-value {
  color: #1f2937;
  font-weight: 600;
}

.role-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.online {
  background: #dcfce7;
  color: #166534;
}

.action-section {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border: none;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  flex: 1;
  min-width: 140px;
  justify-content: center;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.secondary {
  background: #f3f4f6;
  color: #374151;
}

.action-btn.secondary:hover {
  background: #e5e7eb;
}

.action-btn.danger {
  background: #fee2e2;
  color: #dc2626;
}

.action-btn.danger:hover {
  background: #fecaca;
}

.btn-icon {
  font-size: 14px;
}

.stats-card {
  padding: 24px;
}

.stats-card h3 {
  margin: 0 0 20px 0;
  color: #1f2937;
  font-size: 20px;
  font-weight: 700;
}

.stats-grid {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f8fafc;
  border-radius: 12px;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #3b82f6;
  margin-bottom: 8px;
}

.stat-label {
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
}

@media (max-width: 968px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .user-profile-container {
    padding: 16px;
  }
  
  .profile-card {
    padding: 24px;
  }
  
  .avatar-section {
    flex-direction: column;
    text-align: center;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .action-section {
    flex-direction: column;
  }
  
  .action-btn {
    justify-content: center;
  }
}
</style> 