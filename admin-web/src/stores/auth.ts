import { defineStore } from 'pinia'

import { login, type LoginParams } from '../api/auth'

// useAuthStore 管理管理员登录状态。
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('admin_token') || '',
    adminId: Number(localStorage.getItem('admin_id') || 0),
    username: localStorage.getItem('admin_username') || '',
  }),
  getters: {
    // isLoggedIn 判断当前是否已经登录。
    isLoggedIn: (state) => Boolean(state.token),
  },
  actions: {
    // login 调用管理员登录接口并保存登录态。
    async login(params: LoginParams) {
      const result = await login(params)
      this.token = result.token
      this.adminId = result.admin_id
      this.username = result.username
      localStorage.setItem('admin_token', result.token)
      localStorage.setItem('admin_id', String(result.admin_id))
      localStorage.setItem('admin_username', result.username)
    },
    // logout 清理本地登录态。
    logout() {
      this.token = ''
      this.adminId = 0
      this.username = ''
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_id')
      localStorage.removeItem('admin_username')
    },
  },
})
