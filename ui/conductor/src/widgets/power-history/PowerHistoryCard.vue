<template>
  <content-card class="d-flex flex-column pa-8">
    <energy-graph
        v-if="showChart"
        :chart-title="props.chartTitle"
        :demand="chartDemandName"
        :generated="chartGeneratedName"
        class="flex-grow-1"/>
    <power-total
        v-if="showTotal"
        :demand="totalDemandName"
        :generated="totalGeneratedName"
        :occupancy="totalOccupancyName"
        class="mx-auto mt-8"/>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {isNullOrUndef} from '@/util/types.js';
import EnergyGraph from '@/widgets/power-history/PowerHistoryGraph.vue';
import PowerTotal from '@/widgets/power-total/PowerTotal.vue';
import {computed} from 'vue';

const props = defineProps({
  // device name for live electric demand
  demand: {
    type: String,
    default: null
  },
  // device name for demand meter readings, defaults to demand
  demandHistory: {
    type: String,
    default: null
  },
  // device name for live electric generation
  generated: {
    type: String,
    default: null
  },
  // device name for generated/PV meter readings, defaults to generated
  generatedHistory: {
    type: String,
    default: null
  },
  // device name for live occupancy
  occupancy: {
    type: String,
    default: null
  },
  hideChart: {
    type: Boolean,
    default: false
  },
  hideTotal: {
    type: Boolean,
    default: false
  },
  chartTitle: {
    type: String,
    default: undefined
  }
});

const chartGeneratedName = computed(() => props.generatedHistory ?? props.generated);
const chartDemandName = computed(() => props.demandHistory ?? props.demand);

const totalGeneratedName = computed(() => props.generated);
const totalDemandName = computed(() => props.demand);
const totalOccupancyName = computed(() => props.occupancy);

const showChart = computed(() => {
  if (props.hideChart) return false;
  return !isNullOrUndef(chartGeneratedName.value) || !isNullOrUndef(chartDemandName.value);
});
const showTotal = computed(() => {
  if (props.hideTotal) return false;
  return !isNullOrUndef(totalDemandName.value) || !isNullOrUndef(totalGeneratedName.value);
});
</script>

<style scoped>

</style>
