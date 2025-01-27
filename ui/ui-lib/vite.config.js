import vue from '@vitejs/plugin-vue';
import {fileURLToPath, URL} from 'url';
import {defineConfig} from 'vite';
import eslintPlugin from 'vite-plugin-eslint';

// https://vitejs.dev/config/
export default defineConfig(({mode}) => {
  return {
    build: {
      lib: {
        entry: fileURLToPath(new URL('./index.js', import.meta.url)),
        name:
            'sc-bos-ui-lib'
      },
      rollupOptions: {
        external: ['vue'],
        output:
            {
              exports: 'named',
              globals:
                  {
                    vue: 'Vue'
                  }
            }
      }
    },
    plugins: [
      vue(),
      eslintPlugin()
    ]
  };
})
