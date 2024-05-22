<template>
  <span>
    <v-progress-linear indeterminate v-if="loading"/>
    <v-subheader
        v-else-if="!notificationSidebar.length"
        class="text-title-caps-large neutral--text text--lighten-3">
      No Past Notifications
    </v-subheader>
    <div v-else>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">
        Past {{ notificationSidebar.length }} {{
          notificationSidebar.length === 1 ? 'Notification' : 'Notifications'
        }}
      </v-subheader>
      <v-card
          v-for="(data, index) in notificationSidebar"
          :key="index"
          class="mt-4"
          elevation="0">
        <span class="d-flex flex-row flex-nowrap px-4 mb-2">
          <v-icon :class="[data.severity.color, 'mt-n2']" size="22">{{ data.severity.icon }}</v-icon>
          <v-spacer/>
          <v-card-subtitle class="text-caption pa-0 pb-2 grey--text">
            {{ data.created }}
          </v-card-subtitle>
        </span>
        <v-card-subtitle class="ma-0 pa-0 px-4 white--text text-capitalize">
          {{ data.description }}
        </v-card-subtitle>
        <v-divider v-if="index < notificationSidebar.length - 1" class="my-3"/>
      </v-card>
    </div>
  </span>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useSidebarStore} from '@/stores/sidebar';
import {computed} from 'vue';


const sidebar = useSidebarStore();
const notification = useNotifications();

const name = computed(() => sidebar.data?.name);
const item = computed(() => sidebar.data?.item);
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
    const icon = icons[notification.severityData(item.severity).text.toLowerCase()];
    const color = item.resolveTime ?
        'grey--text' :
        notification.severityData(item.severity).color;
    return {
      ...item,
      severity: {icon, color},
      created: timestampToDate(item.createTime).toLocaleString()
    };
  });
});
</script>
