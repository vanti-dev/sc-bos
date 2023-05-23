<template>
  <div id="energy-graph" :style="{width, height}">
    <apexchart type="area" height="100%" :options="options" :series="series" v-if="data.length > 0"/>
    <v-card-text v-else class="error">No data available</v-card-text>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {listMeterReadingHistory} from '@/api/sc/traits/meter-history';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, ref, watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  width: {
    type: String,
    default: '430px'
  },
  height: {
    type: String,
    default: '230px'
  },
  span: { // how wide the bars of the histogram are
    type: Number,
    default: 15 * 60 * 1000 // in ms
  }
});

const pollDelay = computed(() => props.span / 10);
const now = ref(Date.now());
const nowHandle = ref(0);
onMounted(() => {
  nowHandle.value = setInterval(() => {
    now.value = Date.now();
  }, pollDelay.value);
});
onUnmounted(() => {
  clearInterval(nowHandle.value);
});
const baseRequest = computed(() => {
  if (!props.name) return undefined;
  const period = {
    startTime: new Date(now.value - 24 * 60 * 60 * 1000)
  };
  return ({
    name: props.name,
    period,
    pageSize: 1000,
    pageToken: ''
  });
});
const meterHistoryRecords = ref(/** @type {MeterReadingRecord.AsObject[]} */ []);

const pollHandle = ref(0);

async function pollReadings() {
  const req = baseRequest.value;
  const all = [];
  try {
    while (true) {
      const page = await listMeterReadingHistory(req, {});
      all.push(...page.meterReadingRecordsList);
      req.pageToken = page.nextPageToken;
      if (!req.pageToken) {
        break;
      }
    }
  } catch (e) {
    console.error('error getting meter readings', e);
  }
  meterHistoryRecords.value = all;
  pollHandle.value = setTimeout(pollReadings, pollDelay.value);
}

onUnmounted(() => {
  clearTimeout(pollHandle.value);
});

watch(() => baseRequest.value, (baseRequest) => {
  // close existing stream if present
  clearTimeout(pollHandle.value);
  meterHistoryRecords.value = [];

  // create new stream
  if (baseRequest) {
    pollReadings();
  }
}, {immediate: true});

const data = computed(() => {
  const span = props.span;
  const dst = [];
  const records = meterHistoryRecords.value;
  if (records.length > 0) {
    // create a list of data points that show the change in value since the previous reading
    /** @type {MeterReadingRecord.AsObject} */
    let lastReading = null;
    /** @type {MeterReadingRecord.AsObject} */
    let readingCur = null;
    for (const record of records) {
      if (!lastReading) {
        lastReading = record;
        readingCur = record;
        continue;
      }

      // special case if the meter was reset
      if (readingCur.meterReading.usage > record.meterReading.usage) {
        const diff = readingCur.meterReading.usage - lastReading.meterReading.usage;
        dst.push({
          x: new Date(timestampToDate(readingCur.recordTime)),
          y: diff
        });
        lastReading = readingCur = record;
        continue;
      }
      readingCur = record;
      const t0 = timestampToDate(lastReading.recordTime);
      const t1 = timestampToDate(record.recordTime);
      const d = t1 - t0;
      if (d > span) {
        const segmentCount = Math.floor(d / span);
        const diff = (record.meterReading.usage - lastReading.meterReading.usage) / segmentCount;
        lastReading = record;
        dst.push({
          x: new Date(t1),
          y: diff
        });
      }
    }
    // process the last reading, if we haven't already
    const finalReading = records[records.length - 1];
    const t0 = timestampToDate(lastReading.recordTime);
    const t1 = timestampToDate(finalReading.recordTime);
    if (t0 !== t1) {
      const diff = finalReading.meterReading.usage - lastReading.meterReading.usage;
      dst.push({
        x: new Date(t1),
        y: diff
      });
    }
  }
  return dst;
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


// UI Error handling
const errorStore = useErrorStore();
const unwatchErrors = [];
onMounted(() => {
  unwatchErrors.push(errorStore.registerTracker(meterHistoryRecords));
});
onUnmounted(() => {
  for (const unwatchError of unwatchErrors) {
    unwatchError();
  }
});

</script>

<style lang="scss" scoped>
#energy-graph {
  height: 135px;
}
</style>
