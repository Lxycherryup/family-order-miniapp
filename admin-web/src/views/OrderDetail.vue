<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

import { getOrderDetail, updateOrderStatus, type OrderDetail } from '../api/order'
import PageHeader from '../components/PageHeader.vue'

const ORDER_PENDING = 1
const ORDER_COOKING = 2
const ORDER_COMPLETED = 3
const ORDER_CANCELED = 4

const props = defineProps<{
  id: string
}>()

const router = useRouter()
const loading = ref(false)
const order = ref<OrderDetail | null>(null)

const statusText = computed(() => {
  switch (order.value?.status) {
    case ORDER_PENDING:
      return '待处理'
    case ORDER_COOKING:
      return '制作中'
    case ORDER_COMPLETED:
      return '已完成'
    case ORDER_CANCELED:
      return '已取消'
    default:
      return '未知'
  }
})

// loadDetail 加载订单详情。
async function loadDetail() {
  loading.value = true
  try {
    order.value = await getOrderDetail(Number(props.id))
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载订单详情失败')
  } finally {
    loading.value = false
  }
}

// changeStatus 修改订单状态。
async function changeStatus(status: number) {
  if (!order.value) {
    return
  }
  try {
    await updateOrderStatus(order.value.id, status)
    ElMessage.success('订单状态已更新')
    await loadDetail()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '更新订单状态失败')
  }
}

onMounted(loadDetail)
</script>

<template>
  <section v-loading="loading">
    <PageHeader title="订单详情" description="查看订单主信息和菜品明细。">
      <template #actions>
        <el-button @click="router.push('/orders')">返回列表</el-button>
      </template>
    </PageHeader>

    <el-empty v-if="!order" description="订单不存在" />
    <template v-else>
      <section class="detail-panel">
        <div>
          <span>订单号</span>
          <strong>{{ order.order_no }}</strong>
        </div>
        <div>
          <span>用户ID</span>
          <strong>{{ order.user_id }}</strong>
        </div>
        <div>
          <span>订单金额</span>
          <strong>¥{{ Number(order.total_amount).toFixed(2) }}</strong>
        </div>
        <div>
          <span>状态</span>
          <strong>{{ statusText }}</strong>
        </div>
        <div>
          <span>创建时间</span>
          <strong>{{ order.created_at }}</strong>
        </div>
        <div>
          <span>备注</span>
          <strong>{{ order.remark || '无' }}</strong>
        </div>
      </section>

      <section class="toolbar">
        <el-button
          v-if="order.status === ORDER_PENDING"
          type="primary"
          @click="changeStatus(ORDER_COOKING)"
        >
          改为制作中
        </el-button>
        <el-button
          v-if="order.status === ORDER_COOKING"
          type="success"
          @click="changeStatus(ORDER_COMPLETED)"
        >
          改为已完成
        </el-button>
        <el-button
          v-if="order.status === ORDER_PENDING || order.status === ORDER_COOKING"
          type="danger"
          @click="changeStatus(ORDER_CANCELED)"
        >
          取消订单
        </el-button>
      </section>

      <el-table :data="order.items" border class="data-table">
        <el-table-column prop="dish_name" label="菜品" min-width="180" />
        <el-table-column label="单价" width="120">
          <template #default="{ row }">¥{{ Number(row.unit_price).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column label="小计" width="120">
          <template #default="{ row }">¥{{ Number(row.subtotal_amount).toFixed(2) }}</template>
        </el-table-column>
      </el-table>
    </template>
  </section>
</template>
