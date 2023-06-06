<template>
  <v-container fluid class="pa-0 pt-2 d-flex flex-column" style="min-width: 270px; max-width: 1200px">
    <div class="d-flex flex-column flex-md-row">
      <h3 class="text-h3 pt-2 pb-6 flex-grow-1">Status: Building Overview</h3>
      <sc-status-card/>
    </div>
    <div class="d-flex flex-column flex-lg-row">
      <div class="flex-grow-1 d-flex flex-column mr-lg-8">
        <energy-card :generated="supplyZone" :metered="energyZone"/>
      </div>
      <div class="d-flex flex-column" style="min-width: 250px;">
        <occupancy-card :name="occupancyZone"/>
        <environmental-card :name="environmentalZone" :external-name="externalZone"/>
      </div>
    </div>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
import OccupancyCard from '@/routes/ops/components/OccupancyCard.vue';
import EnvironmentalCard from '@/routes/ops/components/EnvironmentalCard.vue';
import EnergyCard from '@/routes/ops/components/EnergyCard.vue';
import ScStatusCard from '@/routes/ops/components/ScStatusCard.vue';
import {useAppConfigStore} from '@/stores/app-config';

const appConfig = useAppConfigStore();

const buildingZone = computed(() => appConfig.config?.ops?.buildingZone ?? '');
const energyZone = buildingZone;
const environmentalZone = buildingZone;
const externalZone = computed(() => environmentalZone.value + '/outside');
const occupancyZone = buildingZone;

const supplyZone = computed(() => energyZone.value + '/supply');

</script>
