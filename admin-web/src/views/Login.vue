<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// form 登录表单数据。
const form = reactive({
  username: 'admin',
  password: '',
})

const submitting = ref(false)

// handleLogin 提交管理员登录。
async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.error('请输入管理员账号和密码')
    return
  }

  submitting.value = true
  try {
    await authStore.login({
      username: form.username,
      password: form.password,
    })
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (err) {
    const message = err instanceof Error ? err.message : '登录失败'
    ElMessage.error(message)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <main class="login-page">
    <section class="login-panel">
      <div class="login-copy">
        <h1>家庭点餐管理后台</h1>
        <p>维护菜品菜单，查看家人订单，管理点餐权限。</p>
      </div>

      <el-form class="login-form" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="管理员账号">
          <el-input v-model="form.username" autocomplete="username" placeholder="请输入管理员账号" />
        </el-form-item>
        <el-form-item label="管理员密码">
          <el-input
            v-model="form.password"
            autocomplete="current-password"
            placeholder="请输入管理员密码"
            show-password
            type="password"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-button
          class="login-button"
          :loading="submitting"
          native-type="submit"
          type="primary"
        >
          登录
        </el-button>
      </el-form>
    </section>
  </main>
</template>
