<template>
  <div class="fill-height layout-overview">
    <header>
      <h3 class="text-h3">Building Status Overview</h3>
    </header>
    <section v-if="showSectionMain" class="section-main">
      <power-history-card
          v-if="powerHistoryConfig"
          style="min-height: 415px;"
          v-bind="powerHistoryConfig"/>
      <occupancy-card
          v-if="occupancyHistoryConfig"
          style="min-height: 415px"
          v-bind="occupancyHistoryConfig"/>
    </section>
    <section v-if="showSectionRight" class="section-right">
      <environmental-card
          v-if="environmentalConfig"
          v-bind="environmentalConfig"
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
  powerHistoryConfig,
  occupancyHistoryConfig,
  environmentalConfig
} = useBuildingConfig();

const showSectionMain = computed(() => Boolean(powerHistoryConfig.value || occupancyHistoryConfig.value));
const showSectionRight = computed(() => Boolean(environmentalConfig.value));
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
