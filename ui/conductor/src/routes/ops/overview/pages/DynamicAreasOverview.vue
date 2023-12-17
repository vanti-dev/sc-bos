<template>
  <v-container fluid class="d-flex flex-column pt-0 pr-0">
    <div class="d-flex flex-row flex-nowrap mb-2">
      <h3 class="text-h3 pt-2 pb-6">
        <span v-for="(value, __, index) of segments" :key="index">
          {{ value }}<span v-if="index < Object.keys(segments).length - 1"> /</span>
        </span>
        Status Overview
      </h3>
    </div>
    <v-row class="ml-0">
      <v-col :class="[{ 'pr-0': !displayRightColumn }, 'ml-0 pl-0']" :style="graphWidth">
        <left-column v-if="displayLeftColumn" :item="overViewStore.getActiveOverview"/>
      </v-col>
      <v-col v-if="displayRightColumn" cols="6" class="mr-0 pr-0 pt-0" style="width: 500px; max-width: 500px;">
        <right-column :item="overViewStore.getActiveOverview"/>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
import {usePageStore} from '@/stores/page';
import {useOverviewStore} from '@/routes/ops/overview/overviewStore';
import LeftColumn from '@/routes/ops/overview/pages/components/LeftColumn.vue';
import RightColumn from '@/routes/ops/overview/pages/components/RightColumn.vue';

const props = defineProps({
  pathSegments: {
    type: Array,
    required: true
  }
});
const overViewStore = useOverviewStore();
const pageStore = usePageStore();
const graphWidth = computed(() => `min-width: calc(100% - 500px - ${pageStore.drawerWidth}px)`);

/**
 * Compute the segments of the path
 *
 * @type {import('vue').ComputedRef<Record<string, string>>}
 */
const segments = computed(() => {
  return props.pathSegments
      .filter(segment => segment !== '') // Exclude empty strings
      .reduce((acc, segment, index) => {
        acc[`prop${index}`] = segment;
        return acc;
      }, {});
});

/**
 * Compute whether to display the left column
 *
 * @type {import('vue').ComputedRef<boolean>} displayLeftColumn
 */
const displayLeftColumn = computed(() => {
  const {getActiveOverview} = overViewStore;
  const emergencyLighting = getActiveOverview?.traits?.showEmergencyLighting;
  const notifications = getActiveOverview?.traits?.showNotifications;
  const lighting = getActiveOverview?.traits?.showLighting;
  const power = getActiveOverview?.traits?.showPower;

  return emergencyLighting || notifications || lighting || power;
});

/**
 * Compute whether to display the right column
 *
 * @type {import('vue').ComputedRef<boolean>} displayRightColumn
 */
const displayRightColumn = computed(() => {
  const {getActiveOverview} = overViewStore;
  const airQuality = getActiveOverview?.traits?.showAirQuality;
  const occupancy = getActiveOverview?.traits?.showOccupancy;
  const energyConsumption = getActiveOverview?.traits?.showEnergyConsumption;
  const environment = getActiveOverview?.traits?.showEnvironment;

  return airQuality || occupancy || energyConsumption || environment;
});
</script>
