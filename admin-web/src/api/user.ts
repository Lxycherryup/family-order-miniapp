import { request } from './request'

export interface User {
  id: number
  openid: string
  nickname: string
  avatar_url: string
  status: number
  is_whitelist: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface UserListParams {
  is_whitelist?: string
  status?: number
  page?: number
  page_size?: number
}

// listUsers 查询小程序用户列表。
export function listUsers(params: UserListParams = {}) {
  return request.get<unknown, User[]>('/users', { params })
}

// updateUserWhitelist 更新用户白名单状态。
export function updateUserWhitelist(id: number, isWhitelist: boolean) {
  return request.put<unknown, boolean>(`/users/${id}/whitelist`, {
    is_whitelist: isWhitelist,
  })
}

// updateUserStatus 更新用户启用状态。
export function updateUserStatus(id: number, status: number) {
  return request.put<unknown, boolean>(`/users/${id}/status`, { status })
}
