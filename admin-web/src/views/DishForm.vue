<script setup lang="ts">
import { reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'

import { createDish, updateDish, type Dish, type DishSaveParams } from '../api/dish'
import type { Category } from '../api/category'

interface Props {
  modelValue: boolean
  dish?: Dish | null
  categories: Category[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const form = reactive<DishSaveParams>({
  category_id: 0,
  name: '',
  description: '',
  image_url: '',
  price: 0,
  unit: '份',
  sort: 100,
})

// resetForm 根据当前菜品重置表单。
function resetForm() {
  form.category_id = props.dish?.category_id || props.categories[0]?.id || 0
  form.name = props.dish?.name || ''
  form.description = props.dish?.description || ''
  form.image_url = props.dish?.image_url || ''
  form.price = props.dish?.price || 0
  form.unit = props.dish?.unit || '份'
  form.sort = props.dish?.sort || 100
}

// closeDialog 关闭菜品表单弹窗。
function closeDialog() {
  emit('update:modelValue', false)
}

// saveDish 保存菜品。
async function saveDish() {
  if (!form.name.trim()) {
    ElMessage.error('请输入菜品名称')
    return
  }
  if (!form.category_id) {
    ElMessage.error('请选择菜品分类')
    return
  }
  if (form.price < 0) {
    ElMessage.error('菜品价格不能小于 0')
    return
  }

  try {
    const payload = {
      ...form,
      name: form.name.trim(),
      unit: form.unit.trim() || '份',
      description: form.description.trim(),
      image_url: form.image_url.trim(),
    }
    if (props.dish) {
      await updateDish(props.dish.id, payload)
    } else {
      await createDish(payload)
    }
    ElMessage.success('菜品已保存')
    emit('saved')
    closeDialog()
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '保存菜品失败')
  }
}

watch(
  () => [props.modelValue, props.dish, props.categories.length],
  () => {
    if (props.modelValue) {
      resetForm()
    }
  },
)
</script>

<template>
  <el-dialog
    :model-value="modelValue"
    :title="dish ? '编辑菜品' : '新增菜品'"
    width="640px"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-form label-position="top">
      <div class="form-grid">
        <el-form-item label="菜品分类" required>
          <el-select v-model="form.category_id" placeholder="请选择分类">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="菜品名称" required>
          <el-input v-model="form.name" maxlength="100" placeholder="请输入菜品名称" />
        </el-form-item>
        <el-form-item label="价格" required>
          <el-input-number v-model="form.price" :min="0" :precision="2" :step="1" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="form.unit" maxlength="16" placeholder="份" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="1" :max="9999" />
        </el-form-item>
        <el-form-item label="图片地址">
          <el-input v-model="form.image_url" maxlength="255" placeholder="https://..." />
        </el-form-item>
      </div>
      <el-form-item label="菜品描述">
        <el-input
          v-model="form.description"
          maxlength="500"
          placeholder="可填写口味、份量或备注"
          :rows="3"
          type="textarea"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="closeDialog">取消</el-button>
      <el-button type="primary" @click="saveDish">保存</el-button>
    </template>
  </el-dialog>
</template>
