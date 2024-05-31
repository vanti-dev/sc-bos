<template>
  <side-bar>
    <v-tabs grow>
      <v-tab>
        <v-icon>mdi-bell</v-icon>
      </v-tab>
      <v-tab>
        <v-icon>mdi-devices</v-icon>
      </v-tab>

      <v-tab-item>
        <past-notifications-tab/>
      </v-tab-item>

      <v-tab-item>
        <device-info-tab :device-id="sidebar.title" :device-data="sidebar.data"/>
      </v-tab-item>
    </v-tabs>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import DeviceInfoTab from '@/routes/ops/notifications/NotificationSideBarTabs/DeviceInfoTab.vue';
import PastNotificationsTab from '@/routes/ops/notifications/NotificationSideBarTabs/PastNotificationsTab.vue';
import {useSidebarStore} from '@/stores/sidebar';
import {usePullMetadata} from '@/traits/metadata/metadata.js';
import {watch} from 'vue';


const sidebar = useSidebarStore();
const {value: metadata} = usePullMetadata(() => sidebar.data?.notification?.item?.source);

watch(metadata, (metadata) => {
  if (metadata) {
    sidebar.data = {
      metadata,
      notification: sidebar.data?.notification
    };
  }
}, {immediate: true, deep: true});
</script>

<style scoped>
:deep(.v-slide-group__prev),
:deep(.v-slide-group__next) {
  display: none !important;
}

</style>
