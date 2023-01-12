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

// Watch for changes in pageStore.sidebarData, which is where the data table item gets passed through to here
watch(sidebarData, async (device) => {
  deviceInfo.value = {};
  if (device && device.hasOwnProperty('metadata')) {
    const data = Object.entries(device.metadata);
    // filter data
    const filtered = data.filter(([key, value]) => {
      // don't display traits or membership
      if (key === 'traitsList' || key === 'membership') {
        return false;
      // ignore empty arrays
      } else if (Array.isArray(value)) {
        return value.length > 0;
      }
      // ignore undefined props
      return value !== undefined;
    });
    // expand and flatten data
    filtered.forEach(([key, value]) => {
      switch (key) {
        case 'location': {
          deviceInfo.value['zone'] = value.title;
          if (value.moreMap.length > 0) {
            for (const more of device.metadata.location.moreMap) {
              deviceInfo.value[more[0]] = more[1];
            }
          }
          break;
        }
        default: {
          deviceInfo.value[key] = value;
        }
      }
    });
  }
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
