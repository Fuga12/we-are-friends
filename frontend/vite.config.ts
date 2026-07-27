import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      // '@' -> папка src. Позволяет писать import '@/components/...'
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    host: true,
    // Все запросы фронта на /api уходят на бэкенд (Litestar, порт 3001).
    // Благодаря этому в разработке нет проблем с CORS, а во фронте можно
    // писать простые относительные пути вида fetch('/api/news').
    proxy: {
      '/api': 'http://localhost:3001',
    },
  },
})
