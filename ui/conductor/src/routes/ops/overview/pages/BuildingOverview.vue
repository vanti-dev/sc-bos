<template>
  <div class="fill-height layout-overview">
    <header>
      <h3 class="text-h3">Building Status Overview</h3>
    </header>
    <section v-if="showSectionMain" class="section-main">
      <power-history-card
          v-if="showEnergy"
          style="min-height: 415px;"
          :demand="energyZone"
          :generated="supplyZone"
          :hide-chart="!showEnergyChart"
          :hide-total="!showEnergyIntensity"/>
      <occupancy-card
          v-if="showOccupancy"
          style="min-height: 415px"
          :source="occupancyZone"/>
    </section>
    <section v-if="showSectionRight" class="section-right">
      <environmental-card
          v-if="showEnvironment"
          :internal="environmentalZone"
          :external="externalZone"
          should-wrap/>
    </section>
  </div>
</template>

<script setup>
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import EnvironmentalCard from '@/widgets/environmental/EnvironmentalCard.vue';
import OccupancyCard from '@/widgets/occupancy/OccupancyCard.vue';
import PowerHistoryCard from '@/widgets/power-history/PowerHistoryCard.vue';
import {computed} from 'vue';

const {
  showEnergy,
  showEnergyChart,
  showEnergyIntensity,
  supplyZone,
  energyZone,
  showOccupancy,
  occupancyZone,
  showEnvironment,
  environmentalZone,
  externalZone
} = useBuildingConfig();

const showSectionMain = computed(() => showEnergy.value || showOccupancy.value);
const showSectionRight = computed(() => showEnvironment.value);

</script>

<style scoped>
.layout-overview {
  display: grid;
  grid-template-columns: 1fr 260px;
  grid-template-rows: auto  1fr;
  gap: 24px;
  width: 100%;
  align-items: stretch;
}

.section-main, .section-right {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-main {
  min-width: 0;
}

.layout-overview > header {
  grid-column: 1 / span 2;
}
</style>
