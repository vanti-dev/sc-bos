<template>
  <v-list tile class="ma-0 pa-0">
    <v-subheader class="text-body-large sidebarTitle pa-4">{{ title }}</v-subheader>
    <slot/>
  </v-list>
</template>

<script setup>
import {defineProps} from 'vue';

defineProps({
  title: {
    type: String,
    default: 'title'
  }
});

</script>

<style scoped>
.sidebarTitle {
  background: var(--v-neutral-lighten1);
  height: auto;
  font-weight: bold;
}
</style>
