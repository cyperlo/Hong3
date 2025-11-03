import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // 监听所有网络接口，允许通过 IP 访问
    port: 5173, // Vite 默认端口
    strictPort: false, // 如果端口被占用，尝试下一个可用端口
  },
})
