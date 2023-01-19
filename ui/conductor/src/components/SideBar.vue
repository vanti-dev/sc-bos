<template>
  <v-list tile class="ma-0 pa-0">
    <v-subheader class="text-body-large sidebar-title pa-4">{{ sidebarTitle }}</v-subheader>
    <slot/>
  </v-list>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';

const pageStore = usePageStore();
const {sidebarTitle} = storeToRefs(pageStore);

</script>

<style scoped>
.sidebar-title {
  background: var(--v-neutral-lighten1);
  height: auto;
  font-weight: bold;
}
</style>
