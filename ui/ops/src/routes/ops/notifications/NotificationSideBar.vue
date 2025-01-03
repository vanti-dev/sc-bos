<template>
  <side-bar>
    <v-tabs grow v-model="tab">
      <v-tab value="notifications">
        <v-icon>mdi-bell</v-icon>
      </v-tab>
      <v-tab value="control">
        <v-icon>mdi-devices</v-icon>
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="notifications">
        <past-notifications-tab/>
      </v-tabs-window-item>
      <v-tabs-window-item value="control">
        <device-info-tab :device-id="sidebar.title" :device-data="sidebar.data"/>
      </v-tabs-window-item>
    </v-tabs-window>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import DeviceInfoTab from '@/routes/ops/notifications/NotificationSideBarTabs/DeviceInfoTab.vue';
import PastNotificationsTab from '@/routes/ops/notifications/NotificationSideBarTabs/PastNotificationsTab.vue';
import {useSidebarStore} from '@/stores/sidebar';
import {usePullMetadata} from '@/traits/metadata/metadata.js';
import {ref, watch} from 'vue';

const tab = ref(null);
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
