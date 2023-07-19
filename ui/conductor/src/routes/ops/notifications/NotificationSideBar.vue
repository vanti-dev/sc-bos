<template>
  <side-bar>
    <v-card v-for="(data, index) in sidebar" :key="index">
      <v-card-title class="text-capitalize">{{ data[0] }}</v-card-title>
      <v-card-subtitle>{{ data[1] ? data[1] : 'N/A' }}</v-card-subtitle>
    </v-card>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';

import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const sidebar = computed(() => {
  return Object.entries(sidebarData.value);
});
</script>

<style scoped>
</style>
