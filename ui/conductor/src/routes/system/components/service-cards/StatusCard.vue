<template>
  <v-card flat tile>
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Status Details</v-subheader>
    <v-card-text class="px-4 pt-3 pb-5">
      <v-row>
        <v-list-item class="py-1">
          <v-list-item-content class="py-0">
            <v-list-item-title class="text-body-small neutral--text text--lighten-4 text-capitalize">
              Status
            </v-list-item-title>
            <v-list-item-subtitle class="text-subtitle-1 py-1 font-weight-medium text-wrap ml-2">
              <service-status :service="sidebarData"/>
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-row>

      <v-row v-for="(value, key) in statusDetails" :key="key">
        <v-list-item class="py-1">
          <v-list-item-content class="py-0">
            <v-list-item-title class="text-body-small neutral--text text--lighten-4 text-capitalize">
              {{ camelToSentence(key) }}
            </v-list-item-title>
            <v-list-item-subtitle class="text-subtitle-1 py-1 font-weight-medium text-wrap ml-2">
              {{ value }}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import ServiceStatus from '@/routes/system/components/ServiceStatus.vue';
import {useSidebarStore} from '@/stores/sidebar';
import {camelToSentence} from '@/util/string';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const sidebar = useSidebarStore();
const {sidebarData} = storeToRefs(sidebar);


const isRunning = computed(() => sidebarData.value.active);
const lastActiveTime = computed(() => timestampToDate(sidebarData.value.lastActiveTime));
const isStopped = computed(() => !isRunning.value);
const lastInactiveTime = computed(() => timestampToDate(sidebarData.value.lastInactiveTime));

const isErrored = computed(() => sidebarData.value.error);


// Computed property for displaying status details
const statusDetails = computed(() => {
  if (isRunning.value) {
    return {
      lastInactiveTime: new Date(lastInactiveTime.value).toLocaleString()
    };
  } else if (isStopped.value) {
    return {
      lastActiveTime: new Date(lastActiveTime.value).toLocaleString()
    };
  } else if (isErrored.value) {
    return {
      lastActiveTime: new Date(lastActiveTime.value).toLocaleString()
    };
  }

  return {
    status: 'Unknown'
  };
});
</script>

<style scoped>
</style>
