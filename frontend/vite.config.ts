import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': process.env.VITE_DEV_API_PROXY_TARGET ?? 'http://localhost:8080',
      '/healthz': process.env.VITE_DEV_API_PROXY_TARGET ?? 'http://localhost:8080',
      '/readyz': process.env.VITE_DEV_API_PROXY_TARGET ?? 'http://localhost:8080',
      '/ws': {
        target: process.env.VITE_DEV_WS_PROXY_TARGET ?? 'ws://localhost:8080',
        ws: true
      }
    }
  }
});
