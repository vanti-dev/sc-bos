<template>
  <v-container fluid class="d-flex flex-row flex-nowrap pl-0 pt-0 pr-0">
    <v-col class="ml-0 pl-0" :style="graphWidth">
      <v-row class="ml-0 pl-0">
        <h3 class="text-h3 pt-2 pb-6">Building Status Overview</h3>
      </v-row>
      <v-row class="d-flex flex-row">
        <v-col cols="12">
          <energy-card
              v-if="showEnergy"
              :show-chart="showEnergyChart"
              :show-intensity="showEnergyIntensity"
              :generated="supplyZone"
              :metered="energyZone"/>
          <occupancy-card v-if="showOccupancy" :name="occupancyZone"/>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="2" class="mr-0 pr-0" style="width: 260px; max-width: 260px;">
      <environmental-card
          v-if="showEnvironment"
          class="ops-sidebar"
          :name="environmentalZone"
          :external-name="externalZone"
          should-wrap/>
    </v-col>
  </v-container>
</template>

<script setup>
import EnergyCard from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyCard.vue';
import EnvironmentalCard from '@/routes/ops/overview/pages/widgets/environmental/EnvironmentalCard.vue';
import OccupancyCard from '@/routes/ops/overview/pages/widgets/occupancy/OccupancyCard.vue';
import {usePageStore} from '@/stores/page';
import {useUiConfigStore} from '@/stores/ui-config';
import {useWidgetsStore} from '@/stores/widgets';
import {computed} from 'vue';

const uiConfig = useUiConfigStore();
const pageStore = usePageStore();
const {activeOverviewWidgets} = useWidgetsStore();
const graphWidth = computed(() => `min-width: calc(100% - 260px - ${pageStore.drawerWidth}px)`);


const buildingZone = computed(() => uiConfig.config?.ops?.buildingZone ?? '');
const showOccupancy = computed(() => activeOverviewWidgets?.showOccupancy);
const occupancyZone = buildingZone;

const showEnergy = computed(() => {
  if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
    return activeOverviewWidgets.showEnergyConsumption;
  } else {
    return activeOverviewWidgets.showEnergyConsumption?.showChart ||
        activeOverviewWidgets.showEnergyConsumption?.showIntensity;
  }
});
const showEnergyChart = computed(() => {
  if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
    return activeOverviewWidgets.showEnergyConsumption;
  } else {
    return activeOverviewWidgets.showEnergyConsumption?.showChart;
  }
});
const showEnergyIntensity = computed(() => {
  if (typeof activeOverviewWidgets.showEnergyConsumption === 'boolean') {
    return activeOverviewWidgets.showEnergyConsumption;
  } else {
    return activeOverviewWidgets.showEnergyConsumption?.showIntensity;
  }
});
const energyZone = buildingZone;
const supplyZone = computed(() => energyZone.value + '/supply');

const showEnvironment = computed(() => activeOverviewWidgets?.showEnvironment);
const environmentalZone = buildingZone;
const externalZone = computed(() => environmentalZone.value + '/outside');

</script>

<style scoped>
.ops-sidebar {
  position: relative;
  top: 63px;
}
</style>
