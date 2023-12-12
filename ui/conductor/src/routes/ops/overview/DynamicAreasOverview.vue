<template>
  <v-container fluid class="d-flex flex-row flex-nowrap pt-0 pr-0">
    <v-col class="ml-0 pl-0" :style="graphWidth">
      <v-row class="ml-0 pl-0">
        <h3 class="text-h3 pt-2 pb-6">
          <span v-for="(value, __, index) of segments" :key="index">
            {{ value }}<span v-if="index < Object.keys(segments).length - 1"> /</span>
          </span>
          Status Overview
        </h3>
      </v-row>
      <v-row class="d-flex flex-row">
        <v-col cols="12"/>
      </v-row>
    </v-col>
    <v-col cols="2" class="mr-0 pr-0 pt-0" style="width: 260px; max-width: 260px;">
      <sc-status-card/>
    </v-col>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';

import {usePageStore} from '@/stores/page';
import ScStatusCard from '@/routes/ops/components/ScStatusCard.vue';

const props = defineProps({
  pathSegments: {
    type: Array,
    required: true
  }
});

const pageStore = usePageStore();
const graphWidth = computed(() => `min-width: calc(100% - 260px - ${pageStore.drawerWidth}px)`);

// Function to convert pathSegments into an object
const segments = computed(() => {
  return props.pathSegments
      .filter(segment => segment !== '') // Exclude empty strings
      .reduce((acc, segment, index) => {
        acc[`prop${index}`] = segment;
        return acc;
      }, {});
});
</script>
