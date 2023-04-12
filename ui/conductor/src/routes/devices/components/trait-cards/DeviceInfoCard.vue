<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" two-line>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Information</v-subheader>
      <v-list-item v-for="(val, key) in deviceInfo" :key="key" class="py-1">
        <v-list-item-content class="py-0">
          <v-list-item-title class="text-capitalize">{{ camelToSentence(key) }}</v-list-item-title>
          <v-list-item-subtitle>{{ val }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {camelToSentence} from '@/util/string';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

// calculate deviceInfo based on sidebarData
const deviceInfo = computed(() => {
  const info = {};
  if (sidebarData?.value?.metadata !== undefined) {
    const data = Object.entries(sidebarData.value.metadata);
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
          info['zone'] = value.title;
          if (value.moreMap.length > 0) {
            for (const more of sidebarData.value.metadata.location.moreMap) {
              info[more[0]] = more[1];
            }
          }
          break;
        }
        default: {
          info[key] = value;
        }
      }
    });
  }
  return info;
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
</style>
