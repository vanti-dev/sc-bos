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
        <device-info-tab :device-id="sidebar.sidebarTitle" :device-data="sidebar.listedDevice"/>
      </v-tab-item>
    </v-tabs>
  </side-bar>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {listDevices} from '@/api/ui/devices';
import SideBar from '@/components/SideBar.vue';
import DeviceInfoTab from '@/routes/ops/notifications/NotificationSideBarTabs/DeviceInfoTab.vue';
import PastNotificationsTab from '@/routes/ops/notifications/NotificationSideBarTabs/PastNotificationsTab.vue';
import {useSidebarStore} from '@/stores/sidebar';
import {reactive, watch} from 'vue';


const sidebar = useSidebarStore();
const listedDevice = reactive(newActionTracker());

watch(() => sidebar.sidebarTitle, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    const newQuery = {
      query: {
        conditionsList: [
          {
            field: 'metadata.name',
            stringEqual: newVal
          }
        ]
      }
    };

    if (newVal) {
      listDevices(newQuery, listedDevice);
    }
  }
}, {immediate: true, deep: true});

watch(() => listedDevice, () => {
  if (listedDevice.response) sidebar.listedDevice = listedDevice?.response?.devicesList[0];
}, {immediate: true, deep: true});
</script>

<style scoped>
:deep(.v-slide-group__prev),
:deep(.v-slide-group__next) {
  display: none !important;
}

</style>
