import { request } from './request'

export interface Order {
  id: number
  order_no: string
  user_id: number
  total_amount: number
  status: number
  remark: string
  created_at: string
  updated_at: string
}

export interface OrderItem {
  id: number
  order_id: number
  dish_id: number
  dish_name: string
  dish_image_url: string
  unit_price: number
  quantity: number
  subtotal_amount: number
  created_at: string
}

export interface OrderDetail extends Order {
  items: OrderItem[]
}

export interface OrderListParams {
  status?: number
  page?: number
  page_size?: number
}

// listOrders 查询订单列表。
export function listOrders(params: OrderListParams = {}) {
  return request.get<unknown, Order[]>('/orders', { params })
}

// getOrderDetail 查询订单详情。
export function getOrderDetail(id: number) {
  return request.get<unknown, OrderDetail>(`/orders/${id}`)
}

// updateOrderStatus 更新订单状态。
export function updateOrderStatus(id: number, status: number) {
  return request.put<unknown, boolean>(`/orders/${id}/status`, { status })
}
