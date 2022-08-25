import vue from '@vitejs/plugin-vue2';
import {VuetifyResolver} from 'unplugin-vue-components/resolvers';
import Components from 'unplugin-vue-components/vite';
import {fileURLToPath, URL} from 'url';
import {defineConfig} from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Components({
      resolvers: [
        VuetifyResolver()
      ]
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  }
});
