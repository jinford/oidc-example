import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 14000,
    proxy: {
      '/api': {
        target: 'http://localhost:14001', // Go サーバーの URL
        changeOrigin: true,
      },
    },
  },
})
