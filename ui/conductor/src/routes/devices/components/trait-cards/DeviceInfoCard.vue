<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0" two-line>
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Information</v-subheader>
      <v-list-item v-for="(val, key) in deviceInfo[0]" :key="key" class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">{{ camelToSentence(key) }}</v-list-item-title>
        <v-list-item-subtitle class="py-1 font-weight-medium">{{ val }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item v-for="(subVal, subKey) in deviceInfo[1]" :key="subKey" class="pt-2">
        <v-list-item-content class="py-0">
          <v-list-item-title class="text-body-small text-capitalize">
            {{ camelToSentence(subKey) }}
          </v-list-item-title>
          <v-list-item-subtitle v-for="(val, key) in subVal" :key="key" class="py-1">
            <v-row class="d-flex flex-row flex-nowrap">
              <v-col class="text-capitalize text-caption" cols="6">
                {{ camelToSentence(key) }}:
              </v-col>
              <v-col class="d-flex flex-column flex-nowrap justify-end ml-n3 pt-4 font-weight-medium text-wrap">
                {{ val }}
              </v-col>
            </v-row>
          </v-list-item-subtitle>
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
  let subInfo = {};
  if (sidebarData?.value?.metadata) {
    const data = Object.entries(sidebarData.value.metadata);
    // filter data
    const filtered = data.filter(([key, value]) => {
      // don't display traits or membership
      if (['traitsList', 'membership'].includes(key)) {
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
      if (key === 'location') {
        info['zone'] = value.title;
        if (value.moreMap.length) {
          for (const more of sidebarData.value.metadata.location.moreMap) {
            info[more[0]] = more[1];
          }
        }
      } else {
        // any inner-object to be flattened
        if (typeof value === 'object') {
          // if there is no sub info yet then add
          if (!subInfo.length) subInfo = {[key]: {}};
          // loop through the object for inner values
          for (const subVal in value) {
            // if we find actual data then collect them
            if (value.hasOwnProperty(subVal) && value[subVal] && subVal !== 'moreMap') {
              subInfo[key] = {...subInfo[key], [subVal]: value[subVal]};
            }
          };
        } else {
          info[key] = value;
        }
      }
    });
  }

  return [info, subInfo];
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
</style>
