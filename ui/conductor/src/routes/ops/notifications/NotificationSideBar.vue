<template>
  <side-bar>
    <div v-if="sidebar.length">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">
        {{ sidebar[0].key }}
      </v-subheader>
      <v-card
          v-for="(data, index) in sidebar[0].value"
          :key="index"
          class="mt-6"
          elevation="0">
        <span class="d-flex flex-row flex-nowrap px-4 mb-4">
          <v-icon :class="[data.severity.color, 'mt-n2']" size="28">{{ icons[data.severity.text] }}</v-icon>
          <v-spacer/>
          <v-card-subtitle class="text-caption pa-0 pb-2 grey--text">
            {{ data.created }}
          </v-card-subtitle>
        </span>
        <v-card-title class="text-subtitle-1 ma-0 pa-0 px-4 mt-2 text-capitalize">
          {{ data.description }}
        </v-card-title>
        <v-divider v-if="index < sidebar[0].value.length - 1" class="mt-4 mb-6"/>
      </v-card>
    </div>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';

import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';


const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);
const notification = useNotifications();

const icons = computed(() => {
  return {
    info: 'mdi-information-outline',
    warn: 'mdi-alert-circle-outline',
    alert: 'mdi-alert-box-outline',
    danger: 'mdi-close-octagon'
  };
});

const sidebar = computed(() => {
  if (!sidebarData || !sidebarData.value || typeof sidebarData.value !== 'object') {
    return []; // or handle the unexpected data structure appropriately
  }

  const filteredEntries = Object.entries(sidebarData.value).filter(([key, value]) => {
    return key === 'past10';
  });

  const mergedData = filteredEntries.map(([key, value]) => {
    const reducedValue = [];

    value.forEach(item => {
      reducedValue.push({
        severity: {
          text: notification.severityData(item.severity).text.toLowerCase(),
          color: notification.severityData(item.severity).color
        },
        description: item.description,
        created: new Date(item.createTime).toLocaleString()
      });
    });

    return {
      key: 'Past 10 notification',
      value: reducedValue
    };
  });

  return mergedData;
});
</script>

<style scoped>
</style>
