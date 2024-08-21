<template>
  <div class="text-body-large font-weight-bold sidebar-title px-4 py-2 bg-neutral-lighten-1 d-flex">
    <span class="sidebar-title-text my-auto">
      {{ sidebar.title }}
    </span>
    <v-spacer/>
    <v-btn icon="mdi-close" variant="text" size="small" @click="sidebar.closeSidebar()" color="neutral-lighten-4">
      <v-icon size="24">mdi-close</v-icon>
    </v-btn>
  </div>
  <v-list tile class="ma-0 pa-0">
    <slot/>
  </v-list>
</template>

<script setup>
import {useSidebarStore} from '@/stores/sidebar';

const sidebar = useSidebarStore();

</script>

<style scoped>
.sidebar-title {
  height: auto;
  align-items: start;
}

.sidebar-title-text {
  display: block;
  overflow-wrap: break-word;
  min-width: 210px;
}
</style>
