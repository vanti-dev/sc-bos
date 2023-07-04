<template>
  <v-container fluid class="pa-0 d-flex flex-column ops-main">
    <!-- Top bar -->
    <div class="d-flex flex-column flex-md-row justify-space-between ma-0 pa-0">
      <h3 class="text-h3 pt-2 pb-6">Status: Building Overview</h3>
      <sc-status-card/>
    </div>
    <v-row class="d-flex flex-row flex-nowrap ops-content ml-0 pl-0">
      <!-- Main contents -->
      <div class="ops-content__main mr-6" :style="contentWidth">
        <energy-card :generated="supplyZone" :metered="energyZone"/>
        <occupancy-card :name="occupancyZone"/>
      </div>
      <div class="ops-sidebar mr-3">
        <environmental-card :name="environmentalZone" :external-name="externalZone"/>
      </div>
    </v-row>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
import {useAppConfigStore} from '@/stores/app-config';
import {usePageStore} from '@/stores/page';

import ScStatusCard from '@/routes/ops/components/ScStatusCard.vue';
import EnvironmentalCard from '@/routes/ops/components/EnvironmentalCard.vue';
import EnergyCard from '@/routes/ops/components/EnergyCard.vue';
import OccupancyCard from '@/routes/ops/components/OccupancyCard.vue';

const appConfig = useAppConfigStore();
const pageStore = usePageStore();

// for more smooth chart resize when expanding the navigation drawer
const contentWidth = computed(() => {
  const drawerWidth = pageStore.drawerWidth;

  return `max-width: calc(82vw - ${drawerWidth}px); transition: max-width .35s ease-in-out .25s; overflow-x: hidden;`;
});

const buildingZone = computed(() => appConfig.config?.ops?.buildingZone ?? '');

const energyZone = buildingZone;
const environmentalZone = buildingZone;
const occupancyZone = buildingZone;

const externalZone = computed(() => environmentalZone.value + '/outside');

const supplyZone = computed(() => energyZone.value + '/supply');

</script>
<style lang="scss">
.ops-main {
  width: 100%;
  .ops-content {
    &__main {
      width: 100%;
      flex-direction: column;
    }
    .ops-sidebar {
      display: block;
      max-width: 250px;
    }
  }
}
</style>
