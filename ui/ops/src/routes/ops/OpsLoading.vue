<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="text-center">
        <v-progress-circular
            v-if="!err"
            indeterminate
            size="64"
            color="primary"/>
        <v-alert v-else density="compact" variant="outlined" color="error">
          {{ err }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import {useEnabledNavItems} from '@/routes/ops/nav.js';
import router from '@/routes/router.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {onMounted, ref} from 'vue';
import {useRoute} from 'vue-router';

const uiConfig = useUiConfigStore();
const route = useRoute();

const err = ref(null);
const enabledNavItems = useEnabledNavItems();

onMounted(() => {
  uiConfig.loadConfig().then(async () => {
    if (route.path === '/ops/loading') {
      let to = null;
      if (uiConfig.pathEnabled('/ops/overview')) {
        to = '/ops/overview';
      } else {
        const page = enabledNavItems.value?.[0];
        if (page) {
          to = page.link;
        }
      }

      if (to) {
        try {
          await router.replace(to);
        } catch {
          // ignore redirects during the routing
        }
      } else {
        err.value = 'You do not have permission to view any operations pages';
      }
    }
  });
});
</script>

<style scoped>

</style>
