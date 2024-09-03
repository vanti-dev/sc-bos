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
import {useControllerStore} from '@/stores/controller';
import {useUiConfigStore} from '@/stores/ui-config';
import {storeToRefs} from 'pinia';
import {onBeforeMount, onMounted} from 'vue';
import {useTheme} from 'vuetify';

const controller = useControllerStore();
const uiConfig = useUiConfigStore();
const {appBranding} = storeToRefs(uiConfig);
const theme = useTheme();

onBeforeMount(async () => {
  await uiConfig.loadConfig();
  const brandColors = appBranding.value?.brandColors;

  // Apply appBranding colors to the vuetify theme
  if (theme && brandColors) {
    // we support two flavours of brandColors:
    // {foo: "#aabbcc"} or {foo: {base: "#aabbcc", lighten1: "#ddeeff"}}
    // The second format is discouraged as it is based on the old Vuetify2 theme format.

    const brandTheme = theme.themes.value['dark'];

    for (const [key, value] of Object.entries(brandColors)) {
      if (typeof value === 'string') {
        brandTheme.colors[key] = /** @type {string} */ value;
      } else {
        for (const [variant, color] of Object.entries(value)) {
          // Vuetify now uses a different variant naming scheme: lighten-1 instead of lighten1.
          // Vuetify also removed the -base suffix from colors.
          // This converts the old format to the new one
          if (variant.startsWith('lighten')) {
            brandTheme.colors[`${key}-lighten-${variant.substring('lighten'.length)}`] = color;
          } else if (variant.includes('darken')) {
            brandTheme.colors[`${key}-darken-${variant.substring('darken'.length)}`] = color;
          } else if (variant === 'base') {
            brandTheme.colors[key] = color;
          } else {
            brandTheme.colors[`${key}-${variant}`] = color;
          }
        }
      }
    }
  }
});
onMounted(() => {
  controller.sync();
});

</script>

<style lang="scss" scoped>
.v-application {
  background: rgb(var(--v-theme-neutral-darken-1));
}
</style>
