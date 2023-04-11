<template>
  <div id="energy-graph">
    <apexchart type="area" height="180" :options="options" :series="series" v-if="data.length > 0"/>
    <v-card-text v-else class="error">No data available</v-card-text>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newActionTracker} from '@/api/resource';
import {listMeterReadingHistory} from '@/api/sc/traits/meter-history';
import {useErrorStore} from '@/components/ui-error/error';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});

const meterHistoryValue = reactive(newActionTracker());

// UI Error handling
const errorStore = useErrorStore();
let unwatchMeterHistoryErrors;
onMounted(() => {
  unwatchMeterHistoryErrors = errorStore.registerTracker(meterHistoryValue);
});
onUnmounted(() => {
  if (unwatchMeterHistoryErrors) unwatchMeterHistoryErrors();
});

watch(() => props.name, (name) => {
  // close existing stream if present
  closeResource(meterHistoryValue);

  // create new stream
  if (name && name !== '') {
    listMeterReadingHistory(name, meterHistoryValue)
        .catch((e) => {});
  }
}, {immediate: true});

const data = computed(() => {
  const res = meterHistoryValue.response;
  if (res?.meterReadingRecordsList) {
    // create a list of data points that show the change in value since the previous reading
    return res.meterReadingRecordsList.map((reading, i) => {
      if (i === 0) {
        return {
          x: new Date(reading.meterReading.endTime.seconds * 1000),
          y: reading.meterReading.usage
        };
      } else {
        return {
          x: new Date(reading.meterReading.endTime.seconds * 1000),
          y: reading.meterReading.usage - res.meterReadingRecordsList[i - 1].meterReading.usage
        };
      }
    });
  } else {
    return [];
  }
});

const options = {
  chart: {
    id: 'energy-chart',
    toolbar: {show: false},
    foreColor: '#fff'
  },
  xaxis: {
    type: 'datetime'
  },
  yaxis: {
    decimalsInFloat: 0
  },
  dataLabels: {enabled: false},
  colors: ['var(--v-primary-base)'],
  stroke: {
    width: 2,
    curve: 'smooth',
    colors: ['var(--v-primary-base)']
  },
  fill: {
    colors: ['var(--v-primary-base)'],
    gradient: {
      enabled: true,
      opacityFrom: 0.55,
      opacityTo: 0.15,
      gradientToColors: ['var(--v-primary-darken1)']
    }
  },
  grid: {
    borderColor: 'var(--v-neutral-lighten2)',
    yaxis: {
      lines: {
        show: false
      }
    },
    padding: {
      top: 0,
      bottom: 25
    }
  },
  tooltip: {
    theme: 'dark'
  }
};

const series = computed(() => {
  return [{
    name: 'Energy Use',
    data: Object.values(data.value)
  }];
});

</script>

<style lang="scss" scoped>
#energy-graph {
  height: 135px;
}
</style>
