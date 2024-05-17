<template>
  <v-list tile class="ma-0 pa-0">
    <v-subheader class="text-body-large font-weight-bold sidebar-title pa-4 neutral lighten-1">
      <span class="sidebar-title-text my-auto">
        {{ sidebar.title }}
      </span>
      <v-spacer/>
      <v-btn icon small @click="sidebar.closeSidebar()" color="neutral lighten-4">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-subheader>
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
