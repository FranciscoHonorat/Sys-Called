import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vitest/config'

const employeesApi = process.env.EMPLOYEES_API_URL ?? 'http://localhost:18081'
const ticketsApi = process.env.TICKETS_API_URL ?? 'http://localhost:18080'

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    proxy: {
      '/api/employees': {
        target: employeesApi,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/employees/, ''),
        cookiePathRewrite: { '/auth': '/api/employees/auth' },
      },
      '/api/tickets': {
        target: ticketsApi,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/tickets/, ''),
      },
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
  },
})
