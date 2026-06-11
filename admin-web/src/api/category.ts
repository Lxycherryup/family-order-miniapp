import { request } from './request'

export interface Category {
  id: number
  name: string
  sort: number
  status: number
  created_at: string
  updated_at: string
}

export interface CategorySaveParams {
  name: string
  sort: number
}

// listCategories 查询分类列表。
export function listCategories() {
  return request.get<unknown, Category[]>('/categories')
}

// createCategory 创建分类。
export function createCategory(params: CategorySaveParams) {
  return request.post<unknown, Category>('/categories', params)
}

// updateCategory 更新分类。
export function updateCategory(id: number, params: CategorySaveParams) {
  return request.put<unknown, Category>(`/categories/${id}`, params)
}

// deleteCategory 删除分类。
export function deleteCategory(id: number) {
  return request.delete<unknown, boolean>(`/categories/${id}`)
}

// updateCategoryStatus 更新分类状态。
export function updateCategoryStatus(id: number, status: number) {
  return request.put<unknown, boolean>(`/categories/${id}/status`, { status })
}
