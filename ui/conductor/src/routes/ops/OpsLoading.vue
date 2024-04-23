<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="text-center">
        <v-progress-circular
            v-if="!err"
            indeterminate
            size="64"
            color="primary"/>
        <v-alert v-else dense outlined color="error">
          {{ err }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import router from '@/routes/router.js';
import {useUiConfigStore} from '@/stores/ui-config.js';
import {onUpdated, ref} from 'vue';
import {useRoute} from 'vue-router/composables';

const uiConfig = useUiConfigStore();
const route = useRoute();

const err = ref(null);

onUpdated(() => {
  uiConfig.loadConfig().then(async () => {
    if (route.path === '/ops/loading') {
      let to = null;
      if (uiConfig.pathEnabled('/ops/overview')) {
        to = '/ops/overview';
      } else if (uiConfig.pathEnabled('/ops/notifications')) {
        to = '/ops/notifications';
      } else if (uiConfig.pathEnabled('/ops/air-quality')) {
        to = '/ops/air-quality';
      } else if (uiConfig.pathEnabled('/ops/emergency-lighting')) {
        to = '/ops/emergency-lighting';
      } else if (uiConfig.pathEnabled('/ops/security')) {
        to = '/ops/security';
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