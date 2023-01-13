<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps neutral--text text--lighten-3">Information</v-subheader>
      <v-list-item v-for="(val, key) in deviceInfo" :key="key" class="py-1">
        <v-list-item-title class="font-weight-bold text-capitalize">{{ camelToSentence(key) }}</v-list-item-title>
        <v-list-item-subtitle>{{ val }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
import {ref, watch} from 'vue';
import {storeToRefs} from 'pinia';
import {usePageStore} from '@/stores/page';
import {camelToSentence} from '@/util/string';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const deviceInfo = ref({});

// Watch for changes in pageStore.sidebarData, which is where the data table item gets passed
watch(sidebarData, (device) => {
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

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
</style>
