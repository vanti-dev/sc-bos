<template>
  <div class="text-body-large font-weight-bold sidebar-title pl-4 pr-2 py-2 bg-neutral-lighten-1 d-flex">
    <span class="sidebar-title-text my-auto mr-auto text-break">
      {{ sidebar.title }}
    </span>
    <slot name="actions"/>
    <v-btn icon="mdi-close" variant="plain" size="small" @click="onCloseClick">
      <v-icon size="24">mdi-close</v-icon>
    </v-btn>
  </div>
  <v-list tile class="ma-0 pa-0">
    <slot/>
  </v-list>
</template>

<script setup>
import {useSidebarStore} from '@/stores/sidebar';

const emit = defineEmits(['close']);
const sidebar = useSidebarStore();

const onCloseClick = () => {
  emit('close');
  sidebar.closeSidebar();
};
</script>

<style scoped>
.sidebar-title {
  height: auto;
  align-items: start;
}

.sidebar-title-text {
  display: block;
  overflow-wrap: break-word;
}
</style>
