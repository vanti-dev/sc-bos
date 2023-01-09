<template>
  <SideBar title="LIT_001">
    <v-subheader class="text-title-caps neutral--text text--lighten-3">Information</v-subheader>
    <v-list-item v-for="(val, key) in deviceInfo" :key="key" class="py-1">
      <v-list-item-title class="font-weight-bold text-capitalize">{{ parseKey(key) }}</v-list-item-title>
      <v-list-item-subtitle>{{ val }}</v-list-item-subtitle>
    </v-list-item>
    <v-subheader class="text-title-caps neutral--text text--lighten-3">State</v-subheader>
    <v-list-item class="py-1">
      <v-list-item-title class="font-weight-bold text-capitalize">Current Temp</v-list-item-title>
      <v-list-item-subtitle>{{ getCurrentTemp(deviceInfo.deviceId) }}</v-list-item-subtitle>
    </v-list-item>
    <v-list-item class="py-1">
      <v-list-item-title class="font-weight-bold text-capitalize">Set Point</v-list-item-title>
      <v-list-item-subtitle>{{ getSetPoint(deviceInfo.deviceId) }}</v-list-item-subtitle>
    </v-list-item>
  </SideBar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {ref, watch} from 'vue';
import {useHvacStore} from '@/routes/devices/hvac/store';
import {storeToRefs} from 'pinia';
import {usePageStore} from '@/stores/page';

const hvacStore = useHvacStore();
const {getSetPoint, getCurrentTemp} = storeToRefs(hvacStore);

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const deviceInfo = ref({});

watch(sidebarData, async (device) => {
  console.log('watch', device.deviceId, hvacStore.getDevice(device.deviceId));
  deviceInfo.value = hvacStore.getDevice(device.deviceId);
});

/**
 * Inserts spaces into camelCased object key
 *
 * @param {string} key
 * @return {string}
 */
function parseKey(key) {
  return key.replace(/([A-Z])/g, ' $1');
}

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
</style>
