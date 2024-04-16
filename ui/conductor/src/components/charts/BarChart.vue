<template>
  <Bar
      :options="props.chartOptions"
      :data="props.chartData"
      :dataset-id-key="props.datasetIdKey"/>
</template>

<script setup>
import {BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip} from 'chart.js';
import {Bar} from 'vue-chartjs';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

const props = defineProps({
  datasetIdKey: {
    type: String,
    default: 'label'
  },
  chartData: {
    type: Object,
    default: () => {
      return {};
    }
  },
  chartOptions: {
    type: Object,
    default: () => {
      return {};
    }
  }
});
</script>
