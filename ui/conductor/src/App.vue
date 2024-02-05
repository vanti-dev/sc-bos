<template>
  <v-app class="app-root">
    <app-bar/>
    <navigation-drawer/>

    <router-view/>
    <error-view/>
  </v-app>
</template>

<script setup>
import AppBar from '@/components/default/AppBar.vue';
import NavigationDrawer from '@/components/default/NavigationDrawer.vue';
import ErrorView from '@/components/ui-error/ErrorView.vue';

import {onMounted, onBeforeMount} from 'vue';
import {storeToRefs} from 'pinia';
import {useAppConfigStore} from '@/stores/app-config';
import {useControllerStore} from '@/stores/controller';
import useVuetify from '@/composables/useVuetify';

const controller = useControllerStore();
const appConfig = useAppConfigStore();
const {appBranding} = storeToRefs(appConfig);
const vuetifyInstance = useVuetify();

onBeforeMount(async () => {
  await appConfig.loadConfig();

  // Access the vuetify instance from the Vue app
  if (vuetifyInstance) {
    const brandColors = appBranding.value?.brandColors;

    if (brandColors) {
      // Loop through each color key (e.g., 'primary', 'secondary', etc.)
      Object.entries(brandColors).forEach(([colorKey, colorVariants]) => {
        // Now loop through each variant of the color (e.g., 'base', 'lighten1', etc.)
        Object.entries(colorVariants).forEach(([variantKey, variantValue]) => {
          // Update the corresponding Vuetify theme color
          // Check if the variantKey exists to safely update the value
          if (vuetifyInstance.theme.themes.dark[colorKey] && variantValue) {
            vuetifyInstance.theme.themes.dark[colorKey][variantKey] = variantValue;
          }
        });
      });
    }
  }
});

onMounted(() => {
  controller.sync();
});

</script>

<style lang="scss" scoped>
.v-application {
  background: var(--v-neutral-darken1);
}
</style>
