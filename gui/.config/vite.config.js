import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from "path";
import tailwind from "tailwindcss"
import autoprefixer from "autoprefixer"

const dirname = __dirname.split("/.config")[0]

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(dirname, "src"),
      '@@': path.resolve(dirname, "wailsjs")
    },
  },
  css: {
    postcss: {
      plugins: [
        tailwind({ config: './.config/tailwind.config.js' }),
        autoprefixer()
      ],
  }
  },
  build: {
    target: 'ESNext',
  },
})
