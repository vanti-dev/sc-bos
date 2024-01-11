<template>
  <v-container fluid class="d-flex flex-row flex-nowrap pt-0 pr-0">
    <v-col class="ml-0 pl-0" :style="graphWidth">
      <v-row class="ml-0 pl-0">
        <h3 class="text-h3 pt-2 pb-6">Building Status Overview</h3>
      </v-row>
      <v-row class="d-flex flex-row">
        <v-col cols="12">
          <energy-card :generated="supplyZone" :metered="energyZone"/>
          <occupancy-card :name="occupancyZone"/>
        </v-col>
      </v-row>
    </v-col>
    <v-col cols="2" class="mr-0 pr-0 pt-0" style="width: 260px; max-width: 260px;">
      <environmental-card class="ops-sidebar mt-5" :name="environmentalZone" :external-name="externalZone" should-wrap/>
    </v-col>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';
import {usePageStore} from '@/stores/page';
import EnvironmentalCard from '@/routes/ops/overview/pages/widgets/environmental/EnvironmentalCard.vue';
import EnergyCard from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyCard.vue';
import OccupancyCard from '@/routes/ops/overview/pages/widgets/occupancy/OccupancyCard.vue';

const appConfig = useAppConfigStore();
const pageStore = usePageStore();
const graphWidth = computed(() => `min-width: calc(100% - 260px - ${pageStore.drawerWidth}px)`);
const buildingZone = computed(() => appConfig.config?.ops?.buildingZone ?? '');
const energyZone = buildingZone;
const environmentalZone = buildingZone;
const occupancyZone = buildingZone;
const externalZone = computed(() => environmentalZone.value + '/outside');
const supplyZone = computed(() => energyZone.value + '/supply');
</script>
