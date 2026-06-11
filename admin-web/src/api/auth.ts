import { request } from './request'

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  admin_id: number
  username: string
}

export interface AdminProfile {
  admin_id: number
  role: string
}

// login 管理员登录。
export function login(params: LoginParams) {
  return request.post<unknown, LoginResult>('/auth/login', params)
}

// getProfile 获取当前管理员资料。
export function getProfile() {
  return request.get<unknown, AdminProfile>('/auth/profile')
}
