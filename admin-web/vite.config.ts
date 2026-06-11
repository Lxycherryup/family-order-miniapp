import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Vite 开发配置，代理后端接口避免本地跨域。
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8000',
        changeOrigin: true,
      },
    },
  },
})
