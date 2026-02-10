import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(({ command }) => ({
  root: __dirname,
  // GitHub Pages serves from /google-interview-prep/ subpath
  // Vercel serves from / (root) — Vercel sets VERCEL=1 env var
  base: process.env.VERCEL ? '/' : '/google-interview-prep/',
  plugins: [react()],
  css: {
    postcss: path.join(__dirname, 'postcss.config.cjs'),
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'markdown': ['react-markdown', 'remark-gfm', 'react-syntax-highlighter'],
        }
      }
    }
  }
}))
