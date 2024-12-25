import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 13000,
    proxy: {
      '/api': {
        target: 'http://localhost:13001', // Go サーバーの URL
        changeOrigin: true,
      },
    },
  },
})
