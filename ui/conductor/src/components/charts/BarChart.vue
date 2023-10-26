<template>
  <Bar
      :css-classes="props.cssClasses"
      :data="props.chartData"
      :dataset-id-key="props.datasetIdKey"
      :options="props.chartOptions"
      :styles="props.styles"/>
</template>

<script setup>
import {Bar} from 'vue-chartjs';
import {Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale} from 'chart.js';
import annotationPlugin from 'chartjs-plugin-annotation';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, annotationPlugin);

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
