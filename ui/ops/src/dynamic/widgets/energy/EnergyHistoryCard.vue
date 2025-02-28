<template>
  <v-card class="px-6 py-4">
    <v-toolbar class="chart-header" color="transparent">
      <slot name="title">
        <v-toolbar-title class="pa-0 mr-auto">{{ props.title }}</v-toolbar-title>
      </slot>
      <v-btn
          icon="mdi-dots-vertical"
          size="small"
          variant="text">
        <v-icon size="24"/>
        <v-menu activator="parent" location="bottom right" offset="8" :close-on-content-click="false">
          <v-card min-width="24em">
            <v-list density="compact">
              <v-list-subheader title="Sources"/>
              <v-list-item
                  v-for="(item, index) in legendItems"
                  :key="index"
                  @click="item.onClick(item.hidden)"
                  :title="item.text">
                <template #prepend>
                  <v-list-item-action start>
                    <v-checkbox-btn :model-value="!item.hidden" readonly :color="item.bgColor" density="compact"/>
                  </v-list-item-action>
                </template>
              </v-list-item>
            </v-list>
            <v-list>
              <v-list-subheader title="Data"/>
              <period-chooser-rows v-model:start="_start" v-model:end="_end" v-model:offset="_offset"/>
            </v-list>
          </v-card>
        </v-menu>
      </v-btn>
    </v-toolbar>
    <v-card-text>
      <div class="chart__container">
        <bar :options="chartOptions" :data="chartData" :plugins="[vueLegendPlugin, themeColorPlugin]"/>
      </div>
    </v-card-text>
    <energy-tooltip :data="tooltipData" :edges="edges" :tick-unit="tickUnit" :unit="unit"/>
  </v-card>
</template>

<script setup>
import {
  computeDatasets,
  useExternalTooltip,
  useThemeColorPlugin,
  useVueLegendPlugin
} from '@/dynamic/widgets/energy/chart.js';
import {useDateScale} from '@/dynamic/widgets/energy/date.js';
import EnergyTooltip from '@/dynamic/widgets/energy/EnergyTooltip.vue';
import PeriodChooserRows from '@/dynamic/widgets/energy/PeriodChooserRows.vue';
import {useDescribeMeterReading} from '@/traits/meter/meter.js';
import {isNullOrUndef} from '@/util/types.js';
import {BarElement, Chart as ChartJS, Legend, LinearScale, TimeScale, Title, Tooltip} from 'chart.js'
import {startOfDay, startOfYear} from 'date-fns';
import {computed, ref, toRef, toValue, watch} from 'vue';
import {Bar} from 'vue-chartjs';
import 'chartjs-adapter-date-fns';
import {useMeterConsumption, useMetersConsumption} from './consumption.js';

ChartJS.register(Title, Tooltip, BarElement, LinearScale, TimeScale, Legend);

const props = defineProps({
  title: {
    type: String,
    default: 'Energy Usage'
  },
  totalConsumptionName: {
    type: String,
    default: undefined,
  },
  totalProductionName: {
    type: String,
    default: undefined,
  },
  subConsumptionNames: {
    type: [Array],
    default: () => [],
  },
  subProductionNames: {
    type: [Array],
    default: () => [],
  },
  start: {
    type: [String, Number, Date],
    default: 'day', // 'month', 'day', etc. meaning 'start of <day>' or a Date-like object
  },
  end: {
    type: [String, Number, Date],
    default: 'day' // 'month', 'day', etc. meaning 'end of <day>' or a Date-like object
  },
  offset: {
    type: [Number, String],
    default: 0, // when start/End is 'month', 'day', etc. offset that value into the past, like 'last month'
  }
});

// we assume here that all the meters share the same unit, so asking about any will be enough.
const nameForDescribe = computed(() => {
  if (!isNullOrUndef(props.totalConsumptionName)) return props.totalConsumptionName;
  if (!isNullOrUndef(props.totalProductionName)) return props.totalProductionName;
  const toName = (item) => {
    if (typeof item === 'string') return item;
    return item.name;
  }
  if (props.subConsumptionNames.length > 0) return toName(props.subConsumptionNames[0]);
  if (props.subProductionNames.length > 0) return toName(props.subProductionNames[0]);
  return undefined;
})
const {response: meterInfo} = useDescribeMeterReading(nameForDescribe);
const unit = computed(() => meterInfo.value?.unit);

/**
 * @template T
 * @param {import('vue').MaybeRefOrGetter<T>} prop
 * @return {import('vue').Ref<T>}
 */
const useLocalProp = (prop) => {
  const local = ref(toValue(prop.value));
  watch(() => toValue(prop), (value) => {
    local.value = value;
  });
  return local;
}
const _start = useLocalProp(toRef(props, 'start'));
const _end = useLocalProp(toRef(props, 'end'));
const _offset = useLocalProp(toRef(props, 'offset'));

const {edges, pastEdges, tickUnit} = useDateScale(_start, _end, _offset);

const totalConsumption = useMeterConsumption(toRef(props, 'totalConsumptionName'), pastEdges);
const totalProduction = useMeterConsumption(toRef(props, 'totalProductionName'), pastEdges);

const subConsumptions = useMetersConsumption(toRef(props, 'subConsumptionNames'), pastEdges);
const subProductions = useMetersConsumption(toRef(props, 'subProductionNames'), pastEdges);

const {
  external: tooltipExternal,
  data: tooltipData,
} = useExternalTooltip(edges, tickUnit, unit);
const {legendItems, vueLegendPlugin} = useVueLegendPlugin();
const {themeColorPlugin} = useThemeColorPlugin();

const chartOptions = computed(() => {
  return /** @type {import('chart.js').ChartOptions} */ {
    responsive: true,
    maintainAspectRatio: false,
    borderRadius: 3,
    borderWidth: 1,
    interaction: {
      mode: 'index', // a single tooltip with all stacked datasets at the same x location in it
    },
    plugins: {
      tooltip: {
        enabled: false,
        external: tooltipExternal
      },
      legend: {
        display: false, // we use a custom legend plugin and vue for this
      }
    },
    scales: {
      y: {
        stacked: true,
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
          offset: false, // bars default to true here, put ticks back inline with grid lines
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
  };
});

const chartData = computed(() => {
  return {
    labels: edges.value,
    datasets: [
      ...computeDatasets('Consumption', totalConsumption, toRef(props, 'subConsumptionNames'), subConsumptions),
      ...computeDatasets('Production', totalProduction, toRef(props, 'subProductionNames'), subProductions, true),
    ]
  };
});
</script>

<style scoped lang="scss">
.chart-header {
  align-items: center;

  :deep(.v-toolbar__content) {
    justify-content: end;
    flex-wrap: wrap;
  }
}

.chart__container {
  min-height: 500px;
  /* The chart seems to have a padding no mater what we do, this gets rid of it */
  margin: -6px;
}
</style>