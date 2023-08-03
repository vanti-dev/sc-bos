<template>
  <v-container fluid class="pa-0">
    <v-row class="px-3 pt-3 mb-5">
      <h3 class="text-h3 py-2">Security Overview</h3>
      <v-spacer/>
      <sc-status-card style="min-width: 248px"/>
    </v-row>
    <content-card class="mb-8 d-flex flex-column pt-6">
      <v-row class="d-flex flex-row align-center mt-0 mb-4 px-6">
        <v-text-field v-model="search" append-icon="mdi-magnify" dense filled hide-details label="Search devices"/>
        <v-spacer/>
        <v-btn-toggle v-model="viewType" dense mandatory>
          <v-btn large text value="list">List View</v-btn>
          <v-btn large text value="map">Map View</v-btn>
        </v-btn-toggle>
        <v-select
            v-model="filterFloor"
            class="ml-4"
            dense
            :disabled="floorList.length <= 1"
            filled
            hide-details
            :items="floorList"
            label="Floor"
            outlined
            style="min-width: 100px; width: 100%; max-width: 170px"/>
        <v-select
            v-model="notificationStateSelection"
            class="ml-4"
            dense
            filled
            hide-details
            :items="['All', 'Alert', 'Offline', 'Open', 'Closed']"
            label="State"
            outlined
            style="max-width: 100px"/>
      </v-row>
      <ListView v-if="viewType === 'list'" :devices="devicesData" :filter="filter"/>
      <MapView v-else/>
    </content-card>
  </v-container>
</template>

<script setup>
import {computed, onMounted, ref} from 'vue';
import ListView from '@/routes/ops/security/components/ListView.vue';
import MapView from '@/routes/ops/security/components/MapView.vue';

import useDevices from '@/composables/useDevices';

import ContentCard from '@/components/ContentCard.vue';
import ScStatusCard from '@/routes/ops/components/ScStatusCard.vue';

const props = defineProps({
  subsystem: {
    type: String,
    default: 'acs'
  },
  filter: {
    type: Function,
    default: () => true
  }
});

const viewType = ref('list');
const notificationStateSelection = ref('All');

const {floorList, filterFloor, search, devicesData} = useDevices(props);

// import Level0 from '@/clients/ew/level0.svg';
</script>

<style scoped></style>
