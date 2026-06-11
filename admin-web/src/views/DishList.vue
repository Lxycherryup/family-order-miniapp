<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import { listCategories, type Category } from '../api/category'
import { deleteDish, listDishes, updateDishStatus, type Dish } from '../api/dish'
import PageHeader from '../components/PageHeader.vue'
import DishForm from './DishForm.vue'

const STATUS_ON = 1
const STATUS_OFF = 2

const loading = ref(false)
const dialogVisible = ref(false)
const editingDish = ref<Dish | null>(null)
const categories = ref<Category[]>([])
const dishes = ref<Dish[]>([])
const filters = reactive({
  category_id: 0,
  status: 0,
})

const categoryNameMap = computed(() => {
  const result = new Map<number, string>()
  for (const category of categories.value) {
    result.set(category.id, category.name)
  }
  return result
})

// loadData 加载分类和菜品列表。
async function loadData() {
  loading.value = true
  try {
    const [categoryList, dishList] = await Promise.all([
      listCategories(),
      listDishes({
        category_id: filters.category_id || undefined,
        status: filters.status || undefined,
      }),
    ])
    categories.value = categoryList
    dishes.value = dishList
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载菜品失败')
  } finally {
    loading.value = false
  }
}

// openCreate 打开新增菜品弹窗。
function openCreate() {
  editingDish.value = null
  dialogVisible.value = true
}

// openEdit 打开编辑菜品弹窗。
function openEdit(row: Dish) {
  editingDish.value = row
  dialogVisible.value = true
}

// toggleStatus 切换菜品上下架状态。
async function toggleStatus(row: Dish) {
  const nextStatus = row.status === STATUS_ON ? STATUS_OFF : STATUS_ON
  try {
    await updateDishStatus(row.id, nextStatus)
    ElMessage.success(nextStatus === STATUS_ON ? '菜品已上架' : '菜品已下架')
    await loadData()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '更新菜品状态失败')
  }
}

// removeDish 下架菜品。
async function removeDish(row: Dish) {
  try {
    await ElMessageBox.confirm(`确认下架菜品「${row.name}」？`, '下架确认', {
      confirmButtonText: '下架',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteDish(row.id)
    ElMessage.success('菜品已下架')
    await loadData()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(err instanceof Error ? err.message : '下架菜品失败')
    }
  }
}

onMounted(loadData)
</script>

<template>
  <section>
    <PageHeader title="菜品管理" description="维护菜品基础信息、价格和上下架状态。">
      <template #actions>
        <el-button type="primary" @click="openCreate">新增菜品</el-button>
      </template>
    </PageHeader>

    <section class="toolbar">
      <el-select v-model="filters.category_id" clearable placeholder="全部分类" @change="loadData">
        <el-option label="全部分类" :value="0" />
        <el-option
          v-for="category in categories"
          :key="category.id"
          :label="category.name"
          :value="category.id"
        />
      </el-select>
      <el-select v-model="filters.status" clearable placeholder="全部状态" @change="loadData">
        <el-option label="全部状态" :value="0" />
        <el-option label="上架" :value="STATUS_ON" />
        <el-option label="下架" :value="STATUS_OFF" />
      </el-select>
      <el-button @click="loadData">刷新</el-button>
    </section>

    <el-table v-loading="loading" :data="dishes" border class="data-table">
      <el-table-column label="菜品" min-width="220">
        <template #default="{ row }">
          <div class="dish-cell">
            <el-image v-if="row.image_url" class="dish-thumb" fit="cover" :src="row.image_url" />
            <div class="dish-thumb placeholder" v-else>无图</div>
            <div>
              <strong>{{ row.name }}</strong>
              <small>{{ row.description || '暂无描述' }}</small>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="分类" min-width="130">
        <template #default="{ row }">
          {{ categoryNameMap.get(row.category_id) || '-' }}
        </template>
      </el-table-column>
      <el-table-column label="价格" width="120">
        <template #default="{ row }">¥{{ Number(row.price).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="unit" label="单位" width="90" />
      <el-table-column prop="sort" label="排序" width="90" />
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="row.status === STATUS_ON ? 'success' : 'info'">
            {{ row.status === STATUS_ON ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" @click="toggleStatus(row)">
            {{ row.status === STATUS_ON ? '下架' : '上架' }}
          </el-button>
          <el-button size="small" type="danger" @click="removeDish(row)">下架</el-button>
        </template>
      </el-table-column>
    </el-table>

    <DishForm
      v-model="dialogVisible"
      :categories="categories"
      :dish="editingDish"
      @saved="loadData"
    />
  </section>
</template>
