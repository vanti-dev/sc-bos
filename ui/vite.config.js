import {fileURLToPath, URL} from 'url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  optimizeDeps: {
    include: [
      "@ew-auth-poc/ui-gen/src/test_grpc_web_pb.js",
      "@ew-auth-poc/ui-gen/src/test_pb.js",
    ]
  },
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  build: {
    commonjsOptions: {
      include: [/@ew-auth-poc\/ui-gen/, /node_modules/]
    }
  }
})
