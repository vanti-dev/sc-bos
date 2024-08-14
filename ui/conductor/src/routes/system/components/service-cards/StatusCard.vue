<template>
  <v-card flat tile>
    <div class="text-subtitle-2 text-title-caps-large text-neutral-lighten-3">Status Details</div>
    <v-card-text class="px-4 pt-3 pb-5">
      <v-row>
        <v-list-item class="py-1">
          <v-list-item-content class="py-0">
            <v-list-item-title class="text-body-small text-neutral-lighten-4 text-capitalize">
              Status
            </v-list-item-title>
            <v-list-item-subtitle class="text-subtitle-1 py-1 font-weight-medium text-wrap ml-2">
              <service-status :service="sidebar.data.service"/>
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </v-row>

      <v-row v-for="(value, key) in statusDetails" :key="key">
        <v-list-item class="py-1">
          <v-list-item-content class="py-0">
            <v-list-item-title class="text-body-small text-neutral-lighten-4 text-capitalize">
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
import {computed} from 'vue';

const sidebar = useSidebarStore();

const service = computed(() => sidebar.data?.service ?? {});
const isRunning = computed(() => service.value.active);
const lastActiveTime = computed(() => timestampToDate(service.value.lastActiveTime));
const isStopped = computed(() => !isRunning.value);
const lastInactiveTime = computed(() => timestampToDate(service.value.lastInactiveTime));
const isErrored = computed(() => service.value.error);


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
