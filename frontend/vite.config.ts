import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: [
      "tracker-vibe.ru",
      "localhost"
    ],
    host: "localhost",
    port: 5000,
  },

  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src")
    }
  }
})
