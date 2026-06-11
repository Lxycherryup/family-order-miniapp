<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import {
  createCategory,
  deleteCategory,
  listCategories,
  updateCategory,
  updateCategoryStatus,
  type Category,
} from '../api/category'
import PageHeader from '../components/PageHeader.vue'

const STATUS_ENABLED = 1
const STATUS_DISABLED = 2

const loading = ref(false)
const dialogVisible = ref(false)
const editingID = ref<number | null>(null)
const categories = ref<Category[]>([])
const form = reactive({
  name: '',
  sort: 100,
})

// loadCategories 加载分类列表。
async function loadCategories() {
  loading.value = true
  try {
    categories.value = await listCategories()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载分类失败')
  } finally {
    loading.value = false
  }
}

// openCreate 打开新增分类弹窗。
function openCreate() {
  editingID.value = null
  form.name = ''
  form.sort = 100
  dialogVisible.value = true
}

// openEdit 打开编辑分类弹窗。
function openEdit(row: Category) {
  editingID.value = row.id
  form.name = row.name
  form.sort = row.sort
  dialogVisible.value = true
}

// saveCategory 保存分类。
async function saveCategory() {
  if (!form.name.trim()) {
    ElMessage.error('请输入分类名称')
    return
  }

  try {
    const payload = {
      name: form.name.trim(),
      sort: Number(form.sort || 100),
    }
    if (editingID.value) {
      await updateCategory(editingID.value, payload)
    } else {
      await createCategory(payload)
    }
    ElMessage.success('分类已保存')
    dialogVisible.value = false
    await loadCategories()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '保存分类失败')
  }
}

// toggleStatus 切换分类启用状态。
async function toggleStatus(row: Category) {
  const nextStatus = row.status === STATUS_ENABLED ? STATUS_DISABLED : STATUS_ENABLED
  try {
    await updateCategoryStatus(row.id, nextStatus)
    ElMessage.success(nextStatus === STATUS_ENABLED ? '分类已启用' : '分类已禁用')
    await loadCategories()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '更新分类状态失败')
  }
}

// removeCategory 删除分类。
async function removeCategory(row: Category) {
  try {
    await ElMessageBox.confirm(`确认删除分类「${row.name}」？`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
    await deleteCategory(row.id)
    ElMessage.success('分类已删除')
    await loadCategories()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error(err instanceof Error ? err.message : '删除分类失败')
    }
  }
}

onMounted(loadCategories)
</script>

<template>
  <section>
    <PageHeader title="分类管理" description="维护菜品分类、展示顺序和启用状态。">
      <template #actions>
        <el-button type="primary" @click="openCreate">新增分类</el-button>
      </template>
    </PageHeader>

    <el-table v-loading="loading" :data="categories" border class="data-table">
      <el-table-column prop="name" label="分类名称" min-width="180" />
      <el-table-column prop="sort" label="排序" width="120" />
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="row.status === STATUS_ENABLED ? 'success' : 'info'">
            {{ row.status === STATUS_ENABLED ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" min-width="180" />
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" @click="toggleStatus(row)">
            {{ row.status === STATUS_ENABLED ? '禁用' : '启用' }}
          </el-button>
          <el-button size="small" type="danger" @click="removeCategory(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingID ? '编辑分类' : '新增分类'" width="420px">
      <el-form label-position="top">
        <el-form-item label="分类名称" required>
          <el-input v-model="form.name" maxlength="64" placeholder="例如：家常菜" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="1" :max="9999" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveCategory">保存</el-button>
      </template>
    </el-dialog>
  </section>
</template>
