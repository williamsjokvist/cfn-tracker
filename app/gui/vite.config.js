import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from "path";
import tailwind from "@tailwindcss/vite"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    tailwind(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(process.cwd(), 'src'),
      '@@': path.resolve(process.cwd(), 'wailsjs'),
      "@runtime": path.resolve(process.cwd(), "wailsjs", "runtime", "runtime.js"),
      "@model": path.resolve(process.cwd(), "wailsjs", "go", "models.ts"),
      "@cmd": path.resolve(process.cwd(), "wailsjs", "go", "cmd"),
    },
  },
  build: {
    target: 'ESNext',
  },
})
