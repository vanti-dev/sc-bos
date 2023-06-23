<template>
  <LineChartGenerator
      :options="props.chartOptions"
      :data="props.chartData"
      :dataset-id-key="props.datasetIdKey"
      :css-classes="props.cssClasses"
      :styles="props.styles"/>
</template>

<script setup>
import {Line as LineChartGenerator} from 'vue-chartjs';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  Filler,
  LinearScale,
  CategoryScale,
  PointElement
} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, LinearScale, Filler, CategoryScale, PointElement);

const props = defineProps({
  datasetIdKey: {
    type: String,
    default: 'label'
  },
  cssClasses: {
    type: String,
    default: 'position-relative'
  },
  styles: {
    type: Object,
    default: () => {
      return {
        height: ''
      };
    }
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
