<template>
  <div style="height: 100%">
    <div class="d-flex flex-row flex-nowrap justify-end align-center mt-3 mb-6">
      <v-card-title class="text-h4 pa-0 mr-auto pl-4">{{ props.chartTitle }}</v-card-title>
      <div v-if="!props.hideLegends" class="legend">
        <v-checkbox
            v-for="(item, index) in legendItems"
            :key="index"
            :input-value="!item.hidden"
            @change="item.onClick"
            :label="item.text"
            :color="item.bgColor"
            hide-details
            class="mt-0"/>
      </div>
      <template v-if="$slots.options">
        <v-divider vertical class="ml-6 mr-2"/>
        <span>
          <slot name="options"/>
        </span>
      </template>
    </div>
    <LineChartGenerator
        :options="props.chartOptions"
        :data="props.chartData"
        :plugins="[htmlLegendPlugin]"
        :dataset-id-key="props.datasetIdKey"
        :css-classes="props.cssClasses"
        :styles="props.styles"/>
  </div>
</template>

<script setup>
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  TimeScale,
  Title,
  Tooltip
} from 'chart.js';
import {computed, ref} from 'vue';
import {Line as LineChartGenerator} from 'vue-chartjs';
import 'chartjs-adapter-date-fns'; // imported for side effects

ChartJS.register(Title, Tooltip, Legend, LineElement, LinearScale, TimeScale, Filler, CategoryScale, PointElement);

const props = defineProps({
  datasetIdKey: {
    type: String,
    default: 'x'
  },
  cssClasses: {
    type: String,
    default: 'position-relative'
  },
  hideLegends: {
    type: Boolean,
    default: false
  },
  styles: {
    type: Object,
    default: () => {
      return {
        height: ''
      };
    }
  },
  chartTitle: {
    type: String,
    default: ''
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

const _chart = ref(null);
const legendItems = computed(() => {
  /** @type {import('chart.js').Chart} */
  const chart = _chart.value;
  if (!chart) return [];
  const items = chart.options.plugins.legend.labels.generateLabels(chart);
  return items.map((item) => {
    return {
      text: item.text,
      hidden: item.hidden,
      bgColor: item.strokeStyle,
      onClick: (e) => {
        const {type} = chart.config;
        if (type === 'pie' || type === 'doughnut') {
          // Pie and doughnut charts only have a single dataset and visibility is per item
          chart.setDatasetVisibility(item.index, e);
        } else {
          chart.setDatasetVisibility(item.datasetIndex, e);
        }
        chart.update();
      }
    };
  });
});

const htmlLegendPlugin = {
  id: 'htmlLegend',
  afterUpdate(chart) {
    _chart.value = chart;
  }
};
</script>

<style scoped>
.legend {
  display: inline-flex;
  gap: 24px;
}
</style>
