<template>
  <div class="chart__container">
    <line-chart :data="chartData" :options="chartOptions" :plugins="[themeColorPlugin]"/>
  </div>
</template>

<script setup>
import {useDateScale} from '@/components/charts/date.js';
import {useThemeColorPlugin} from '@/components/charts/plugins.js';
import {defineChartOptions} from '@/components/charts/util.js';
import {useAirQualityHistoryMetric} from '@/dynamic/widgets/environmental/airQuality.js';
import {useLocalProp} from '@/util/vue.js';
import {sentenceCase} from 'change-case';
import {Chart as ChartJS, Legend, LinearScale, LineElement, PointElement, TimeScale, Title, Tooltip} from 'chart.js';
import {startOfDay, startOfYear} from 'date-fns';
import {computed, toRef} from 'vue';
import {Line as LineChart} from 'vue-chartjs';
import 'chartjs-adapter-date-fns'; // imported for side effects

ChartJS.register(Title, Tooltip, LineElement, LinearScale, PointElement, TimeScale, Legend);

const props = defineProps({
  source: {
    type: String,
    default: null
  },
  metric: {
    type: String,
    default: 'score'
  },
  unit: {
    type: String,
    default: null, // default calculated based on metric
  },
  start: {
    type: [String, Number, Date],
    default: 'day', // 'month', 'day', etc. meaning 'start of <day>' or a Date-like object
  },
  end: {
    type: [String, Number, Date],
    default: 'day', // 'month', 'day', etc. meaning 'end of <day>' or a Date-like object
  },
  offset: {
    type: [String, Number],
    default: 0, // when start/End is 'month', 'day', etc. offset that value into the past, like 'last month'
  }
});

const _start = useLocalProp(toRef(props, 'start'));
const _end = useLocalProp(toRef(props, 'end'));
const _offset = useLocalProp(toRef(props, 'offset'));

const {edges, pastEdges, tickUnit} = useDateScale(_start, _end, _offset);

const sourceMetrics = useAirQualityHistoryMetric(toRef(props, 'source'), toRef(props, 'metric'), pastEdges);

const yAxisLabel = computed(() => {
  const s = sentenceCase(props.metric);
  if (props.unit) {
    return `${s} (${props.unit})`;
  }
  return s;
});

const {themeColorPlugin} = useThemeColorPlugin();
const chartOptions = computed(() => {
  return defineChartOptions({
    responsive: true,
    maintainAspectRatio: false,
    borderRadius: 3,
    borderWidth: 1,
    interaction: {
      mode: 'index', // a single tooltip with all stacked datasets at the same x location in it
    },
    plugins: {
      legend: {
        display: false, // we use a custom legend plugin and vue for this
      }
    },
    scales: {
      y: {
        stacked: true,
        beginAtZero: true,
        title: {
          display: true,
          text: yAxisLabel.value
        },
        border: {
          color: 'transparent'
        },
        grid: {
          color(ctx) {
            if (ctx.tick.value === 0) return '#fff4';
            return '#fff1';
          },
          drawTicks: false,
        },
        ticks: {
          callback(value) {
            return new Intl.NumberFormat(undefined, {}).format(Math.abs(value));
          },
          color: '#fff',
          padding: 8
        },
      },
      x: {
        type: 'time',
        stacked: true,
        grid: {
          color: '#fff1'
        },
        ticks: {
          maxTicksLimit: 11,
          includeBounds: true,
          callback(value) {
            const unit = tickUnit.value;
            if (unit === 'month' && value === startOfYear(value).getTime()) return this.format(value, this.options.time.displayFormats['year']);
            if (unit === 'hour' && value === startOfDay(value).getTime()) return this.format(value, this.options.time.displayFormats['day']);
            return this.format(value);
          },
          color: '#fff',
          padding: 8,
          maxRotation: 0
        },
        time: {
          unit: tickUnit.value,
          displayFormats: {
            hour: 'H:mm', // default: 4:00 AM
            day: 'd MMM', // default: Feb 10
            month: 'MMM', // default: Feb 2025 - we fix the ambiguity in ticks.callback
          }
        }
      }
    }
  });
});
const chartLabels = computed(() => edges.value.slice(0, -1));
const chartData = computed(() => {
  return {
    labels: chartLabels.value,
    datasets: [
      {
        label: 'Air Quality',
        data: sourceMetrics.value.map(data => data.y)
      }
    ]
  }
})
</script>

<style scoped>

</style>