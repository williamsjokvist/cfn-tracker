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
      '@': '/src',
      '@@': '/wailsjs',
      "@runtime": path.join("/", "wailsjs", "runtime", "runtime.js"),
      "@model": path.join("/", "wailsjs", "go", "models.ts"),
      "@cmd": path.join("/", "wailsjs", "go", "cmd"),
    },
  },
  build: {
    target: 'ESNext',
  },
})
