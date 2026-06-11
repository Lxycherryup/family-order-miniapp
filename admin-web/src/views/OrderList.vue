<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { listOrders, updateOrderStatus, type Order } from '../api/order'
import PageHeader from '../components/PageHeader.vue'

const ORDER_PENDING = 1
const ORDER_COOKING = 2
const ORDER_COMPLETED = 3
const ORDER_CANCELED = 4

const router = useRouter()
const loading = ref(false)
const orders = ref<Order[]>([])
const filters = reactive({
  status: 0,
  page: 1,
  page_size: 20,
})

const statusMap: Record<number, { label: string; type: 'warning' | 'primary' | 'success' | 'info' }> = {
  [ORDER_PENDING]: { label: '待处理', type: 'warning' },
  [ORDER_COOKING]: { label: '制作中', type: 'primary' },
  [ORDER_COMPLETED]: { label: '已完成', type: 'success' },
  [ORDER_CANCELED]: { label: '已取消', type: 'info' },
}

// loadOrders 加载订单列表。
async function loadOrders() {
  loading.value = true
  try {
    orders.value = await listOrders({
      status: filters.status || undefined,
      page: filters.page,
      page_size: filters.page_size,
    })
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载订单失败')
  } finally {
    loading.value = false
  }
}

// changeStatus 修改订单状态。
async function changeStatus(row: Order, status: number) {
  try {
    await updateOrderStatus(row.id, status)
    ElMessage.success('订单状态已更新')
    await loadOrders()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '更新订单状态失败')
  }
}

// showDetail 打开订单详情。
function showDetail(row: Order) {
  router.push(`/orders/${row.id}`)
}

onMounted(loadOrders)
</script>

<template>
  <section>
    <PageHeader title="订单管理" description="查看订单并处理待制作、完成和取消状态。" />

    <section class="toolbar">
      <el-select v-model="filters.status" clearable placeholder="全部状态" @change="loadOrders">
        <el-option label="全部状态" :value="0" />
        <el-option label="待处理" :value="ORDER_PENDING" />
        <el-option label="制作中" :value="ORDER_COOKING" />
        <el-option label="已完成" :value="ORDER_COMPLETED" />
        <el-option label="已取消" :value="ORDER_CANCELED" />
      </el-select>
      <el-button @click="loadOrders">刷新</el-button>
    </section>

    <el-table v-loading="loading" :data="orders" border class="data-table">
      <el-table-column prop="order_no" label="订单号" min-width="190" />
      <el-table-column prop="user_id" label="用户ID" width="100" />
      <el-table-column label="金额" width="120">
        <template #default="{ row }">¥{{ Number(row.total_amount).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="statusMap[row.status]?.type || 'info'">
            {{ statusMap[row.status]?.label || '未知' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="remark" label="备注" min-width="160" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" min-width="180" />
      <el-table-column label="操作" width="340" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="showDetail(row)">详情</el-button>
          <el-button
            v-if="row.status === ORDER_PENDING"
            size="small"
            type="primary"
            @click="changeStatus(row, ORDER_COOKING)"
          >
            制作中
          </el-button>
          <el-button
            v-if="row.status === ORDER_COOKING"
            size="small"
            type="success"
            @click="changeStatus(row, ORDER_COMPLETED)"
          >
            完成
          </el-button>
          <el-button
            v-if="row.status === ORDER_PENDING || row.status === ORDER_COOKING"
            size="small"
            type="danger"
            @click="changeStatus(row, ORDER_CANCELED)"
          >
            取消
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </section>
</template>
