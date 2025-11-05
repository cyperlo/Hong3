import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // 监听所有网络接口，允许通过 IP 访问
    port: 5173, // Vite 默认端口
    strictPort: false, // 如果端口被占用，尝试下一个可用端口
    proxy: {
      // 代理 API 请求到后端
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      // 代理 WebSocket 请求到后端
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
        changeOrigin: true,
      },
    },
  },
})
