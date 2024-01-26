import vue from '@vitejs/plugin-vue2';
import glob from 'glob';
import {createRequire} from 'module';
import {dirname, join, relative} from 'path';
import {VuetifyResolver} from 'unplugin-vue-components/resolvers';
import Components from 'unplugin-vue-components/vite';
import {fileURLToPath, URL} from 'url';
import {defineConfig} from 'vite';
import eslintPlugin from 'vite-plugin-eslint';
import gitRevision from 'vite-plugin-git-revision';

const _require = createRequire(import.meta.url);

// Any import that resolves to a local filesystem dependency that isn't an ESM dependency needs to be in this list.
const optimizeDepsInclude = [];
// Typically that includes local proto files that are referenced via either `file:../` dependencies in package.json
// or via yarn/npm linking the generated sources into this project (or both).
// This snippet will find all *_pb.js files and ensure that they will be handled correctly by vite.
for (const dep of ['@sc-bos/ui-gen', '@smart-core-os/sc-api-grpc-web']) {
  // find proto files in projects
  const protoDirRoot = dirname(_require.resolve(dep + '/package.json'));
  const protoFiles = glob.sync(join(protoDirRoot, '!(node_modules)/**/*_pb.js'))
      .map(p => dep + '/' + relative(protoDirRoot, p));
  optimizeDepsInclude.push(...protoFiles);
  // remove the .js extension so import statements without .js still use the bundle
  optimizeDepsInclude.push(...protoFiles.map(f => f.substring(0, f.length - 3)));
}

// https://vitejs.dev/config/
export default defineConfig({
  define: {},
  optimizeDeps: {
    include: optimizeDepsInclude
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/sass/variables.scss";\n`
      },
      sass: {
        additionalData: `@import "@/sass/variables.scss"\n`
      }
    }
  },
  plugins: [
    vue(),
    eslintPlugin(),
    gitRevision({
      versionCommand: 'describe --tags --always --match ui/*'
    }),
    // can't fix imported var names, so tell eslint to ignore them
    // eslint-disable-next-line new-cap
    Components({
      // eslint-disable-next-line new-cap
      resolvers: [VuetifyResolver()]
    })
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
});
