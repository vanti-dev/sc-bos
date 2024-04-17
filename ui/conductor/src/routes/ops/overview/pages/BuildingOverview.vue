<template>
  <div class="fill-height layout-overview">
    <header>
      <h3 class="text-h3">Building Status Overview</h3>
    </header>
    <section v-if="showSectionMain" class="section-main">
      <power-history-card
          v-if="showEnergy"
          style="min-height: 415px;"
          :show-chart="showEnergyChart"
          :show-intensity="showEnergyIntensity"
          :generated="supplyZone"
          :metered="energyZone"/>
      <occupancy-card
          v-if="showOccupancy"
          style="min-height: 415px"
          :name="occupancyZone"/>
    </section>
    <section v-if="showSectionRight" class="section-right">
      <environmental-card
          v-if="showEnvironment"
          :name="environmentalZone"
          :external-name="externalZone"
          should-wrap/>
    </section>
  </div>
</template>

<script setup>
import {useUiConfigStore} from '@/stores/ui-config';
import {useWidgetsStore} from '@/stores/widgets';
import EnvironmentalCard from '@/widgets/environmental/EnvironmentalCard.vue';
import OccupancyCard from '@/widgets/occupancy/OccupancyCard.vue';
import PowerHistoryCard from '@/widgets/power-history/PowerHistoryCard.vue';
import {computed} from 'vue';

const uiConfig = useUiConfigStore();
const {activeOverviewWidgets} = useWidgetsStore();

const buildingZoneSource = computed(() => uiConfig.config?.ops?.buildingZone ?? '');
const supplyZoneSource = computed(() => uiConfig.config?.ops?.supplyZone ?? '');

const buildingZone = computed(() => buildingZoneSource.value);
const supplyZone = computed(() => supplyZoneSource.value ? supplyZoneSource.value + '/supply' : '');
const energyZone = buildingZone;
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

const showOccupancy = computed(() => activeOverviewWidgets?.showOccupancy);
const occupancyZone = buildingZone;

const showEnvironment = computed(() => activeOverviewWidgets?.showEnvironment);
const environmentalZone = buildingZone;
const externalZone = computed(() => environmentalZone.value + '/outside');

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
