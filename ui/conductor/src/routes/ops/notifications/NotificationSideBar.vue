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
        <PastNotificationsTab/>
      </v-tab-item>

      <v-tab-item>
        <DeviceInfoTab :device-id="pageStore.sidebarTitle" :device-data="pageStore.listedDevice"/>
      </v-tab-item>
    </v-tabs>
  </side-bar>
</template>

<script setup>
import {reactive, watch} from 'vue';
import {listDevices} from '@/api/ui/devices';
import {newActionTracker} from '@/api/resource';
import {usePageStore} from '@/stores/page';


import SideBar from '@/components/SideBar.vue';
import DeviceInfoTab from '@/routes/ops/notifications/NotificationSideBarTabs/DeviceInfoTab.vue';
import PastNotificationsTab from '@/routes/ops/notifications/NotificationSideBarTabs/PastNotificationsTab.vue';


const pageStore = usePageStore();
const listedDevice = reactive(newActionTracker());

watch(() => pageStore.sidebarTitle, (newVal, oldVal) => {
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
  if (listedDevice.response) pageStore.listedDevice = listedDevice?.response?.devicesList[0];
}, {immediate: true, deep: true});
</script>

<style scoped>
:deep(.v-slide-group__prev),
:deep(.v-slide-group__next) {
  display: none !important;
}

</style>
