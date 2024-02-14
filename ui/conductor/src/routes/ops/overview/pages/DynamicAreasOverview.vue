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
import {usePageStore} from '@/stores/page';
import {useUiConfigStore} from '@/stores/ui-config';
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

const uiConfig = useUiConfigStore();
const overviewChildren = computed(() => uiConfig.config?.ops?.overview?.children);

const pageStore = usePageStore();
const graphWidth = computed(() => `min-width: calc(100% - 500px - ${pageStore.drawerWidth}px)`);

/**
 * Compute whether to display the left column
 *
 * @type {import('vue').ComputedRef<boolean>} displayLeftColumn
 */
const displayLeftColumn = computed(() => {
  const emergencyLighting = activeOverview.value?.widgets?.showEmergencyLighting;
  const notifications = activeOverview.value?.widgets?.showNotifications;
  const lighting = activeOverview.value?.widgets?.showLighting;
  const power = activeOverview.value?.widgets?.showPower;

  return emergencyLighting || notifications || lighting || power;
});

/**
 * Compute whether to display the right column
 *
 * @type {import('vue').ComputedRef<boolean>} displayRightColumn
 */
const displayRightColumn = computed(() => {
  const airQuality = activeOverview.value?.widgets?.showAirQuality;
  const occupancy = activeOverview.value?.widgets?.showOccupancy;
  const energyConsumption = activeOverview.value?.widgets?.showEnergyConsumption;
  const environment = activeOverview.value?.widgets?.showEnvironment;

  return airQuality || occupancy || energyConsumption || environment;
});

const findActiveOverview = computed(() => {
  if (!props.pathSegments.length || !overviewChildren.value.length) return null;

  const encodePathSegments = props.pathSegments.map(encodeURIComponent);

  return findActiveItem(overviewChildren.value, encodePathSegments);
});

watch(findActiveOverview, (overviewChild) => {
  // Set the active overview or null if not found
  activeOverview.value = overviewChild;
}, {immediate: true, deep: true, flush: 'sync'});
</script>
