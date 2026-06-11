<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const menus = computed(() => [
  { path: '/categories', label: '分类管理' },
  { path: '/dishes', label: '菜品管理' },
  { path: '/orders', label: '订单管理' },
  { path: '/users', label: '白名单' },
])

// handleLogout 退出当前管理员登录。
async function handleLogout() {
  authStore.logout()
  await router.replace('/login')
}
</script>

<template>
  <main class="dashboard-page">
    <aside class="sidebar">
      <div class="brand">
        <span class="brand-mark">家</span>
        <div>
          <strong>家庭点餐</strong>
          <small>管理后台</small>
        </div>
      </div>
      <nav class="nav-list">
        <router-link
          v-for="menu in menus"
          :key="menu.path"
          class="nav-item"
          :to="menu.path"
        >
          {{ menu.label }}
        </router-link>
      </nav>
    </aside>

    <section class="workspace">
      <header class="topbar">
        <div>
          <h1>家庭点餐管理后台</h1>
          <p>维护菜单、处理订单并管理家人点餐权限。</p>
        </div>
        <div class="admin-actions">
          <span>{{ authStore.username }}</span>
          <el-button @click="handleLogout">退出</el-button>
        </div>
      </header>

      <router-view />
    </section>
  </main>
</template>
