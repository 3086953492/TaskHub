import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import { useAuth } from './composables/useAuth'

const app = createApp(App)

// 初始化认证状态
const { initAuth } = useAuth()
initAuth()

app.use(router).mount('#app')
