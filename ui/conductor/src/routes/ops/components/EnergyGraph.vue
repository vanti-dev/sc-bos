<template>
  <div id="energy-graph">
    <apexchart type="area" height="180" :options="options" :series="series"/>
  </div>
</template>

<script setup>
import {computed, reactive, ref, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {listMeterReadingHistory} from '@/api/sc/traits/meter-history';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  }
});

const meterHistoryValue = reactive(newResourceValue());

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(meterHistoryValue);

  // create new stream
  if (name && name !== '') {
    const res = await listMeterReadingHistory(name, meterHistoryValue);
    // create a list of data points that show the change in value since the previous reading
    data.value = res.meterReadingRecordsList.map((reading, i) => {
      if (i === 0) {
        return {
          x: new Date(reading.meterReading.endTime.seconds*1000),
          y: reading.meterReading.usage
        };
      } else {
        return {
          x: new Date(reading.meterReading.endTime.seconds*1000),
          y: reading.meterReading.usage - res.meterReadingRecordsList[i - 1].meterReading.usage
        };
      }
    });
  }
}, {immediate: true});

const data = ref([]);

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
