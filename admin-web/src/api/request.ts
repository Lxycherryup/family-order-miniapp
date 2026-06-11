import axios from 'axios'

// request 管理端接口请求实例。
export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/admin',
  timeout: 10000,
})

// 请求拦截器，自动携带管理员 Token。
request.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器，统一处理后端响应结构。
request.interceptors.response.use((response) => {
  const body = response.data
  if (body.code !== 0) {
    return Promise.reject(new Error(body.msg || '请求失败'))
  }
  return body.data
})
