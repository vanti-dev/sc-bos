<template>
  <v-container fluid class="d-flex flex-column pt-0 pr-0">
    <div class="d-flex flex-row flex-nowrap mb-2">
      <h3 class="text-h3 pt-2 pb-6">
        {{ overViewStore.getActiveOverview?.title }} Status Overview
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
import {useOverviewStore} from '@/routes/ops/overview/overviewStore';
import LeftColumn from '@/routes/ops/overview/pages/components/LeftColumn.vue';
import RightColumn from '@/routes/ops/overview/pages/components/RightColumn.vue';
import {useAppConfigStore} from '@/stores/app-config';
import {usePageStore} from '@/stores/page';
import {findActiveItem} from '@/util/router';
import {storeToRefs} from 'pinia';
import {computed, watch} from 'vue';

const props = defineProps({
  pathSegments: {
    type: Array,
    required: true
  }
});
const overViewStore = useOverviewStore();
const {activeOverview} = storeToRefs(overViewStore);

const appConfigStore = useAppConfigStore();
const buildingChildren = computed(() => appConfigStore.config?.building?.children || []);

const pageStore = usePageStore();
const graphWidth = computed(() => `min-width: calc(100% - 500px - ${pageStore.drawerWidth}px)`);

/**
 * Compute whether to display the left column
 *
 * @type {import('vue').ComputedRef<boolean>} displayLeftColumn
 */
const displayLeftColumn = computed(() => {
  const emergencyLighting = activeOverview.value?.traits?.showEmergencyLighting;
  const notifications = activeOverview.value?.traits?.showNotifications;
  const lighting = activeOverview.value?.traits?.showLighting;
  const power = activeOverview.value?.traits?.showPower;

  return emergencyLighting || notifications || lighting || power;
});

/**
 * Compute whether to display the right column
 *
 * @type {import('vue').ComputedRef<boolean>} displayRightColumn
 */
const displayRightColumn = computed(() => {
  const airQuality = activeOverview.value?.traits?.showAirQuality;
  const occupancy = activeOverview.value?.traits?.showOccupancy;
  const energyConsumption = activeOverview.value?.traits?.showEnergyConsumption;
  const environment = activeOverview.value?.traits?.showEnvironment;

  return airQuality || occupancy || energyConsumption || environment;
});

const findActiveOverview = computed(() => {
  return findActiveItem(buildingChildren.value, props.pathSegments);
});

watch(() => props.pathSegments, () => {
  // If no path segments are provided, we are on the root overview page
  if (props.pathSegments.length === 0) {
    activeOverview.value = null;
  } else {
    activeOverview.value = findActiveOverview.value;
  }
}, {immediate: true, deep: true});
</script>
