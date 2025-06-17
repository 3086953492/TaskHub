<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import NavBar from './components/NavBar.vue'
import { useAuth } from './composables/useAuth'

const route = useRoute()
const { isLoggedIn } = useAuth()

// 不需要显示导航栏的页面
const hideNavbarRoutes = ['/login', '/register']
const shouldShowNavbar = computed(() => {
  return isLoggedIn.value && !hideNavbarRoutes.includes(route.path)
})
</script>

<template>
  <div id="app">
    <!-- 在需要时显示导航栏 -->
    <NavBar v-if="shouldShowNavbar" />
    
    <!-- 路由视图 -->
    <router-view />
  </div>
</template>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  width: 100%;
  margin: 0;
  padding: 0;
  background: #f8fafc;
}

* {
  box-sizing: border-box;
}

html, body {
  width: 100%;
  margin: 0;
  padding: 0;
  background: #f8fafc;
}
</style>
