<template>
  <v-container fluid class="pa-0 pt-2 d-flex flex-column" style="min-width: 270px; max-width: 1200px">
    <div class="d-flex flex-column flex-md-row">
      <h3 class="text-h3 pt-2 pb-6 flex-grow-1">Status: Building Overview</h3>
      <sc-status-card/>
    </div>
    <div class="d-flex flex-column flex-lg-row">
      <div class="flex-grow-1 d-flex flex-column mr-lg-8">
        <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
          <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">System Monitor</h4>
          <building-status/>
        </content-card>
        <energy-card :current-energy="currentEnergy"/>
      </div>
      <div class="d-flex flex-column" style="min-width: 250px;">
        <occupancy-card :name="zoneName"/>
        <environmental-card :name="zoneName"/>
      </div>
    </div>
    <!-- todo: remove test slider, hook up real backend -->
    <div>
      <span class="text-title-caps">Test</span>
      <v-slider v-model="sliderVal" max="1" step="0.01"/>
    </div>
  </v-container>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {computed, ref} from 'vue';
import BuildingStatus from '@/clients/ew/BuildingStatus_EW.vue';
import OccupancyCard from '@/routes/ops/components/OccupancyCard.vue';
import EnvironmentalCard from '@/routes/ops/components/EnvironmentalCard.vue';
import EnergyCard from '@/routes/ops/components/EnergyCard.vue';
import ScStatusCard from '@/routes/ops/components/ScStatusCard.vue';

const sliderVal = ref(0);

const currentEnergy = computed(() => sliderVal.value*120);

// todo: should probably get this from somewhere
const zoneName = 'Floor 3'; // building-wide zone to use for collecting data

</script>
