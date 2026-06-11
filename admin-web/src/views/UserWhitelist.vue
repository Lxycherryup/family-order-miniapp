<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

import {
  listUsers,
  updateUserStatus,
  updateUserWhitelist,
  type User,
} from '../api/user'
import PageHeader from '../components/PageHeader.vue'

const STATUS_ENABLED = 1
const STATUS_DISABLED = 2

const loading = ref(false)
const users = ref<User[]>([])
const filters = reactive({
  is_whitelist: '',
  status: 0,
  page: 1,
  page_size: 20,
})

// maskOpenID 脱敏展示 OpenID。
function maskOpenID(openid: string) {
  if (openid.length <= 12) {
    return openid
  }
  return `${openid.slice(0, 6)}****${openid.slice(-6)}`
}

// loadUsers 加载小程序用户列表。
async function loadUsers() {
  loading.value = true
  try {
    users.value = await listUsers({
      is_whitelist: filters.is_whitelist || undefined,
      status: filters.status || undefined,
      page: filters.page,
      page_size: filters.page_size,
    })
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载用户失败')
  } finally {
    loading.value = false
  }
}

// changeWhitelist 更新用户白名单。
async function changeWhitelist(row: User, value: boolean) {
  try {
    await updateUserWhitelist(row.id, value)
    row.is_whitelist = value
    ElMessage.success(value ? '已加入白名单' : '已移出白名单')
  } catch (err) {
    row.is_whitelist = !value
    ElMessage.error(err instanceof Error ? err.message : '更新白名单失败')
  }
}

// changeStatus 更新用户状态。
async function changeStatus(row: User) {
  const nextStatus = row.status === STATUS_ENABLED ? STATUS_DISABLED : STATUS_ENABLED
  try {
    await updateUserStatus(row.id, nextStatus)
    row.status = nextStatus
    ElMessage.success(nextStatus === STATUS_ENABLED ? '用户已启用' : '用户已禁用')
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '更新用户状态失败')
  }
}

onMounted(loadUsers)
</script>

<template>
  <section>
    <PageHeader title="白名单管理" description="管理小程序用户点餐权限和账号启用状态。" />

    <section class="toolbar">
      <el-select v-model="filters.is_whitelist" clearable placeholder="全部白名单" @change="loadUsers">
        <el-option label="全部白名单" value="" />
        <el-option label="白名单用户" value="true" />
        <el-option label="非白名单用户" value="false" />
      </el-select>
      <el-select v-model="filters.status" clearable placeholder="全部状态" @change="loadUsers">
        <el-option label="全部状态" :value="0" />
        <el-option label="启用" :value="STATUS_ENABLED" />
        <el-option label="禁用" :value="STATUS_DISABLED" />
      </el-select>
      <el-button @click="loadUsers">刷新</el-button>
    </section>

    <el-table v-loading="loading" :data="users" border class="data-table">
      <el-table-column label="OpenID" min-width="180">
        <template #default="{ row }">{{ maskOpenID(row.openid) }}</template>
      </el-table-column>
      <el-table-column prop="nickname" label="昵称" min-width="140" />
      <el-table-column label="白名单" width="130">
        <template #default="{ row }">
          <el-switch
            v-model="row.is_whitelist"
            inline-prompt
            active-text="是"
            inactive-text="否"
            @change="(value: string | number | boolean) => changeWhitelist(row, Boolean(value))"
          />
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="row.status === STATUS_ENABLED ? 'success' : 'info'">
            {{ row.status === STATUS_ENABLED ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="最近登录" min-width="180" />
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="changeStatus(row)">
            {{ row.status === STATUS_ENABLED ? '禁用' : '启用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>
