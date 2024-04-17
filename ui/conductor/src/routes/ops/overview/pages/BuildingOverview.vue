<template>
  <layout-main-side v-bind="widgets"/>
</template>

<script setup>
import LayoutMainSide from '@/layout/LayoutMainSide.vue';
import useBuildingConfig from '@/routes/ops/overview/pages/buildingConfig.js';
import EnvironmentalCard from '@/widgets/environmental/EnvironmentalCard.vue';
import OccupancyCard from '@/widgets/occupancy/OccupancyCard.vue';
import PowerHistoryCard from '@/widgets/power-history/PowerHistoryCard.vue';
import {computed} from 'vue';

const {
  title,
  powerHistoryConfig,
  occupancyHistoryConfig,
  environmentalConfig
} = useBuildingConfig();

const widgets = computed(() => ({
  title: title.value,
  mainWidgetMinHeight: 415,
  sideWidth: 260,
  main: [
    {component: PowerHistoryCard, props: powerHistoryConfig.value},
    {component: OccupancyCard, props: occupancyHistoryConfig.value}
  ],
  after: [
    {component: EnvironmentalCard, props: environmentalConfig.value}
  ]
}));
</script>

<style scoped>
</style>
