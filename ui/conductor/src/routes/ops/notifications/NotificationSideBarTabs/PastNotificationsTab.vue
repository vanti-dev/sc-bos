<template>
  <div>
    <v-progress-linear indeterminate v-if="loading"/>
    <div
        v-else-if="!hasSource"
        class="text-subtitle-2 text-title-caps-large text-neutral-lighten-3">
      No Notification Selected
    </div>
    <div
        v-else-if="!notificationSidebar.length"
        class="text-subtitle-2 text-title-caps-large text-neutral-lighten-3">
      No Past Notifications
    </div>
    <div v-else>
      <div class="text-subtitle-2 text-title-caps-large text-neutral-lighten-3 pa-4">
        Past {{ notificationSidebar.length }} {{
          notificationSidebar.length === 1 ? 'Notification' : 'Notifications'
        }}
      </div>
      <v-card
          v-for="(data, index) in notificationSidebar"
          :key="index"
          class="mt-4"
          elevation="0">
        <span class="d-flex flex-row align-center flex-nowrap px-4 mb-2">
          <v-icon :class="data.severity.color" size="22">{{ data.severity.icon }}</v-icon>
          <v-spacer/>
          <v-card-subtitle class="text-caption pa-0 text-grey">
            {{ data.created }}
          </v-card-subtitle>
        </span>
        <v-card-text class="ma-0 pa-0 px-4 text-white text-capitalize">
          {{ data.description }}
        </v-card-text>
        <v-divider v-if="index < notificationSidebar.length - 1" class="my-3"/>
      </v-card>
    </div>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {severityData} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useSidebarStore} from '@/stores/sidebar';
import {computed} from 'vue';


const sidebar = useSidebarStore();

const name = computed(() => sidebar.data?.notification?.name);
const item = computed(() => sidebar.data?.notification?.item);
const hasSource = computed(() => Boolean(item.value?.source));
// todo: don't fetch data if we don't have a source
const query = computed(() => ({source: item.value?.source}));
const {pageItems, pageSize, targetItemCount, loading} = useAlertsApi(name, query);
pageSize.value = 10;
targetItemCount.value = 10;


const icons = {
  info: 'mdi-information-outline',
  warn: 'mdi-alert-circle-outline',
  alert: 'mdi-alert-box-outline',
  danger: 'mdi-close-octagon'
};
const notificationSidebar = computed(() => {
  if (pageItems.value.length === 0) return [];

  return pageItems.value.map(item => {
    const icon = icons[severityData(item.severity).text.toLowerCase()];
    const color = item.resolveTime ?
        'text-grey' :
        severityData(item.severity).color;
    return {
      ...item,
      severity: {icon, color},
      created: timestampToDate(item.createTime).toLocaleString()
    };
  });
});
</script>
