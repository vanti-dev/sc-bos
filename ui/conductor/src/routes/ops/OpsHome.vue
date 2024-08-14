<template>
  <router-view/>
</template>

<script setup>
import router from '@/routes/router.js';
import {useUiConfigStore} from '@/stores/ui-config.js';
import {onUpdated} from 'vue';
import {useRoute} from 'vue-router';

const uiConfig = useUiConfigStore();
const route = useRoute();
onUpdated(() => {
  uiConfig.loadConfig().then(async () => {
    if (route.path === '/ops/overview') {
      const pages = uiConfig.getOrDefault('ops.pages');
      let to;
      if (!pages) {
        to = '/ops/overview/building';
      } else {
        to = `/ops/overview/${encodeURIComponent(pages[0].path ?? pages[0].title)}`;
      }
      try {
        await router.replace(to);
      } catch {
        // ignore redirects during the routing
      }
    }
  });
});
</script>
