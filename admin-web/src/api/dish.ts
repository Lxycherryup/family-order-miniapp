import { request } from './request'

export interface Dish {
  id: number
  category_id: number
  name: string
  description: string
  image_url: string
  price: number
  unit: string
  status: number
  sort: number
  created_at: string
  updated_at: string
}

export interface DishSaveParams {
  category_id: number
  name: string
  description: string
  image_url: string
  price: number
  unit: string
  sort: number
}

export interface DishListParams {
  category_id?: number
  status?: number
}

// listDishes 查询菜品列表。
export function listDishes(params: DishListParams = {}) {
  return request.get<unknown, Dish[]>('/dishes', { params })
}

// getDish 查询菜品详情。
export function getDish(id: number) {
  return request.get<unknown, Dish>(`/dishes/${id}`)
}

// createDish 创建菜品。
export function createDish(params: DishSaveParams) {
  return request.post<unknown, Dish>('/dishes', params)
}

// updateDish 更新菜品。
export function updateDish(id: number, params: DishSaveParams) {
  return request.put<unknown, Dish>(`/dishes/${id}`, params)
}

// deleteDish 下架菜品。
export function deleteDish(id: number) {
  return request.delete<unknown, boolean>(`/dishes/${id}`)
}

// updateDishStatus 更新菜品状态。
export function updateDishStatus(id: number, status: number) {
  return request.put<unknown, boolean>(`/dishes/${id}/status`, { status })
}
