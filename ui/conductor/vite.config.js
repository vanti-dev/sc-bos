import {URL, fileURLToPath} from 'url';
import {defineConfig} from 'vite';
import vue from '@vitejs/plugin-vue2';
import Components from 'unplugin-vue-components/vite';
import {VuetifyResolver} from 'unplugin-vue-components/resolvers';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Components({
      // unplugin-vue-components hasn't get updated to support @vitejs/plugin-vue2, it looks for vite-plugin-vue2 to
      // set this for you.
      transformer: 'vue2',
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
