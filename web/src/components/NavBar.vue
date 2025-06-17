<template>
    <nav class="navbar">
      <div class="nav-container">
        <!-- 品牌Logo -->
        <div class="nav-brand">
          <router-link to="/" class="brand-link">
            <h2>TaskHub</h2>
          </router-link>
        </div>
  
        <!-- 导航链接 -->
        <div class="nav-links">
          <router-link 
            to="/" 
            class="nav-link"
            :class="{ active: $route.path === '/' }"
          >
            <span class="nav-icon">📋</span>
            <span>任务列表</span>
          </router-link>
          
          <router-link 
            to="/profile" 
            class="nav-link"
            :class="{ active: $route.path === '/profile' }"
          >
            <span class="nav-icon">👤</span>
            <span>用户信息</span>
          </router-link>
        </div>
  
        <!-- 用户信息和退出 -->
        <div class="nav-user">
          <div v-if="user" class="user-info">
            <span class="user-name">{{ user.nickname || user.username }}</span>
            <span class="user-role" :class="roleClass">{{ roleText }}</span>
          </div>
          
          <button @click="handleLogout" class="logout-btn">
            <span class="logout-icon">🚪</span>
            <span>退出</span>
          </button>
        </div>
  
        <!-- 移动端菜单按钮 -->
        <button @click="toggleMobileMenu" class="mobile-menu-btn">
          <span class="hamburger" :class="{ active: showMobileMenu }"></span>
        </button>
      </div>
  
      <!-- 移动端菜单 -->
      <div class="mobile-menu" :class="{ active: showMobileMenu }">
        <router-link 
          to="/" 
          class="mobile-nav-link"
          @click="closeMobileMenu"
        >
          <span class="nav-icon">📋</span>
          <span>任务列表</span>
        </router-link>
        
        <router-link 
          to="/profile" 
          class="mobile-nav-link"
          @click="closeMobileMenu"
        >
          <span class="nav-icon">👤</span>
          <span>用户信息</span>
        </router-link>
  
        <div class="mobile-user-info">
          <div v-if="user" class="user-info">
            <span class="user-name">{{ user.nickname || user.username }}</span>
            <span class="user-role" :class="roleClass">{{ roleText }}</span>
          </div>
          
          <button @click="handleLogout" class="mobile-logout-btn">
            <span class="logout-icon">🚪</span>
            <span>退出登录</span>
          </button>
        </div>
      </div>
    </nav>
  </template>
  
  <script setup lang="ts">
  import { ref, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import { useAuth } from '../composables/useAuth'
  
  const router = useRouter()
  const { user, logout } = useAuth()
  
  // 移动端菜单状态
  const showMobileMenu = ref(false)
  
  // 用户角色显示
  const roleText = computed(() => {
    return user.value?.role === 'admin' ? '管理员' : '普通用户'
  })
  
  const roleClass = computed(() => {
    return user.value?.role === 'admin' ? 'role-admin' : 'role-user'
  })
  
  // 切换移动端菜单
  const toggleMobileMenu = () => {
    showMobileMenu.value = !showMobileMenu.value
  }
  
  // 关闭移动端菜单
  const closeMobileMenu = () => {
    showMobileMenu.value = false
  }
  
  // 处理退出登录
  const handleLogout = () => {
    logout()
    closeMobileMenu()
    router.push('/login')
  }
  </script>
  
  <style scoped>
  .navbar {
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    border-bottom: 1px solid #e5e7eb;
    position: sticky;
    top: 0;
    z-index: 100;
  }
  
  .nav-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 20px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
  }
  
  .nav-brand .brand-link {
    text-decoration: none;
    color: inherit;
  }
  
  .nav-brand h2 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea, #764ba2);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }
  
  .nav-links {
    display: flex;
    gap: 8px;
  }
  
  .nav-link {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 8px;
    text-decoration: none;
    color: #6b7280;
    transition: all 0.2s;
    font-weight: 500;
  }
  
  .nav-link:hover {
    background: #f3f4f6;
    color: #1f2937;
  }
  
  .nav-link.active {
    background: #e0e7ff;
    color: #3730a3;
  }
  
  .nav-icon {
    font-size: 16px;
  }
  
  .nav-user {
    display: flex;
    align-items: center;
    gap: 16px;
  }
  
  .user-info {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 2px;
  }
  
  .user-name {
    font-weight: 600;
    color: #1f2937;
    font-size: 14px;
  }
  
  .user-role {
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
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
  
  .logout-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background: #f3f4f6;
    border: none;
    border-radius: 8px;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 14px;
  }
  
  .logout-btn:hover {
    background: #e5e7eb;
    color: #374151;
  }
  
  .logout-icon {
    font-size: 14px;
  }
  
  .mobile-menu-btn {
    display: none;
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
  }
  
  .hamburger {
    display: block;
    width: 24px;
    height: 2px;
    background: #374151;
    position: relative;
    transition: all 0.3s;
  }
  
  .hamburger::before,
  .hamburger::after {
    content: '';
    position: absolute;
    width: 24px;
    height: 2px;
    background: #374151;
    transition: all 0.3s;
  }
  
  .hamburger::before {
    top: -8px;
  }
  
  .hamburger::after {
    bottom: -8px;
  }
  
  .hamburger.active {
    background: transparent;
  }
  
  .hamburger.active::before {
    transform: rotate(45deg);
    top: 0;
  }
  
  .hamburger.active::after {
    transform: rotate(-45deg);
    bottom: 0;
  }
  
  .mobile-menu {
    display: none;
    background: white;
    border-top: 1px solid #e5e7eb;
    padding: 16px 20px;
  }
  
  .mobile-menu.active {
    display: block;
  }
  
  .mobile-nav-link {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 0;
    text-decoration: none;
    color: #374151;
    border-bottom: 1px solid #f3f4f6;
    font-weight: 500;
  }
  
  .mobile-nav-link:last-of-type {
    border-bottom: none;
  }
  
  .mobile-user-info {
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #e5e7eb;
  }
  
  .mobile-user-info .user-info {
    align-items: flex-start;
    margin-bottom: 12px;
  }
  
  .mobile-logout-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 12px;
    background: #f3f4f6;
    border: none;
    border-radius: 8px;
    color: #6b7280;
    cursor: pointer;
    font-size: 14px;
    transition: background-color 0.2s;
  }
  
  .mobile-logout-btn:hover {
    background: #e5e7eb;
  }
  
  @media (max-width: 768px) {
    .nav-links,
    .nav-user {
      display: none;
    }
    
    .mobile-menu-btn {
      display: block;
    }
    
    .nav-container {
      padding: 0 16px;
    }
  }
  </style>