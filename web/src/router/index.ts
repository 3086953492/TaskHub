import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import UserProfile from '../views/UserProfile.vue'
import EditProfile from '../views/EditProfile.vue'
import { useAuth } from '../composables/useAuth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'UserProfile',
    component: UserProfile,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/edit',
    name: 'EditProfile',
    component: EditProfile,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const { isLoggedIn } = useAuth()
  
  // 如果路由需要认证且用户未登录，重定向到登录页
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    next('/login')
  }
  // 如果用户已登录且访问登录/注册页，重定向到首页
  else if (isLoggedIn.value && (to.name === 'Login' || to.name === 'Register')) {
    next('/')
  }
  else {
    next()
  }
})

export default router 