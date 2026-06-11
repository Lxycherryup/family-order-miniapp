<script setup lang="ts">
import { useRouter } from 'vue-router'

import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

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
        <a class="nav-item active">工作台</a>
        <a class="nav-item">分类管理</a>
        <a class="nav-item">菜品管理</a>
        <a class="nav-item">订单管理</a>
        <a class="nav-item">白名单</a>
      </nav>
    </aside>

    <section class="workspace">
      <header class="topbar">
        <div>
          <h1>工作台</h1>
          <p>当前后端接口和登录态已经接入，下一步实现业务管理页面。</p>
        </div>
        <div class="admin-actions">
          <span>{{ authStore.username }}</span>
          <el-button @click="handleLogout">退出</el-button>
        </div>
      </header>

      <section class="quick-grid">
        <article class="quick-card">
          <strong>分类管理</strong>
          <p>维护菜品分类、排序和启用状态。</p>
        </article>
        <article class="quick-card">
          <strong>菜品管理</strong>
          <p>维护菜品图片、价格、单位和上下架状态。</p>
        </article>
        <article class="quick-card">
          <strong>订单管理</strong>
          <p>查看订单详情并处理制作状态。</p>
        </article>
      </section>
    </section>
  </main>
</template>
