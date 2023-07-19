<template>
  <side-bar>
    <div v-if="sidebar.length">
      <v-card-title
          class="text-h6 grey--text font-weight-bold
           ma-0 mt-4 pa-0 pl-4 text-capitalize">
        {{ sidebar[0].key }}
      </v-card-title>
      <v-card
          v-for="(data, index) in sidebar[0].value"
          :key="index"
          class="mt-4"
          elevation="0">
        <span v-for="(alert, alertIndex) in Object.entries(data)" :key="alertIndex">
          <v-card-title class="ma-0 pa-0 pl-4 text-capitalize">
            {{ alert[0] }}
          </v-card-title>
          <v-card-text class="grey--text">
            {{ alert[1] }}
          </v-card-text>
        </span>
        <v-divider v-if="index < sidebar[0].value.length - 1" class="mt-2 mb-6"/>
      </v-card>
    </div>
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
  if (!sidebarData || !sidebarData.value || typeof sidebarData.value !== 'object') {
    return []; // or handle the unexpected data structure appropriately
  }

  const filteredEntries = Object.entries(sidebarData.value).filter(([key, value]) => {
    return key === 'past10';
  });

  const mergedData = filteredEntries.map(([key, value]) => {
    return {
      key: 'Past 10 notification',
      value
    };
  });

  return mergedData;
});
</script>

<style scoped>
</style>
