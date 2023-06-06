<template>
  <div id="energy-graph" :style="{width, height}">
    <apexchart
        type="area"
        height="100%"
        :options="options"
        :series="series"/>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {listMeterReadingHistory} from '@/api/sc/traits/meter-history';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, ref, watch} from 'vue';

const props = defineProps({
  name: {
    type: [Array, String],
    default: ''
  },
  width: {
    type: String,
    default: '430px'
  },
  height: {
    type: String,
    default: '275px'
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

// Return an array of object with request details
// This loops through the props.name array if its array
// or only creates an object if its a string
const baseRequest = computed(() => {
  if (!props.name || !props.name.length) return undefined;

  const period = {
    startTime: new Date(now.value - 24 * 60 * 60 * 1000)
  };
  const req = [];

  if (Array.isArray(props.name)) {
    props.name.forEach(name => {
      req.push({
        name,
        period,
        pageSize: 1000,
        pageToken: ''
      });
    });
  } else {
    req.push({
      name: props.name,
      period,
      pageSize: 1000,
      pageToken: ''
    });
  }

  return req;
});
const meterHistoryRecords = ref(/** @type {MeterReadingRecord.AsObject[]} */ []);
const supplyHistoryRecords = ref(/** @type {MeterReadingRecord.AsObject[]} */ []);

const meterPollHandle = ref(0);
const supplyPollHandle = ref(0);
const message = ref('');

/**
 *
 * @param {*} req
 * @param {string} type
 */
async function pollReadings(req, type) {
  const generated = [];
  const supplied = [];
  try {
    while (true) {
      const page = await listMeterReadingHistory(req, {});

      if (type === 'supply') supplied.push(...page.meterReadingRecordsList);
      else generated.push(...page.meterReadingRecordsList);
      req.pageToken = page.nextPageToken;
      if (!req.pageToken) {
        break;
      }
    }
  } catch (e) {
    console.error('error getting meter readings', e);
  }

  if (type === 'supply') {
    supplyHistoryRecords.value = supplied;
    supplyPollHandle.value = setTimeout(pollReadings, pollDelay.value);
  } else {
    meterHistoryRecords.value = generated;
    meterPollHandle.value = setTimeout(pollReadings, pollDelay.value);
  }

  message.value = 'No data available';
}

onUnmounted(() => {
  clearTimeout(supplyPollHandle.value);
  clearTimeout(meterPollHandle.value);
});

watch(() => baseRequest.value, (baseRequest) => {
  baseRequest.forEach(request => {
    // close existing stream if present
    message.value = 'Pulling data';
    meterHistoryRecords.value = [];
    clearTimeout(meterPollHandle.value);

    supplyHistoryRecords.value = [];
    clearTimeout(supplyPollHandle.value);

    // create new stream
    if (request) {
      if (request.name.includes('supply')) pollReadings(request, 'supply');
      else pollReadings(request);
    }
  });
}, {immediate: true, deep: true, flush: 'sync'});

//
//
// Generated energy
const meterData = computed(() => {
  const span = props.span;
  const dst = [];
  const meterRecords = meterHistoryRecords.value;

  if (meterRecords.length > 0) {
    // create a list of data points that show the change in value since the previous reading
    /** @type {MeterReadingRecord.AsObject} */
    let lastReading = null;
    /** @type {MeterReadingRecord.AsObject} */
    let readingCur = null;
    for (const record of meterRecords) {
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
    const finalReading = meterRecords[meterRecords.length - 1];
    const t0 = timestampToDate(lastReading.recordTime);
    const t1 = timestampToDate(finalReading.recordTime);
    if (t0.getTime() !== t1.getTime()) {
      const diff = finalReading.meterReading.usage - lastReading.meterReading.usage;
      dst.push({
        x: new Date(t1),
        y: diff
      });
    }
  }
  return dst;
});

//
//
// PV energy
const supplyData = computed(() => {
  const span = props.span;
  const dst = [];
  const supplyRecords = supplyHistoryRecords.value;

  if (supplyRecords.length > 0) {
    // create a list of data points that show the change in value since the previous reading
    /** @type {MeterReadingRecord.AsObject} */
    let lastReading = null;
    /** @type {MeterReadingRecord.AsObject} */
    let readingCur = null;
    for (const record of supplyRecords) {
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
    const finalReading = supplyRecords[supplyRecords.length - 1];
    const t0 = timestampToDate(lastReading.recordTime);
    const t1 = timestampToDate(finalReading.recordTime);
    if (t0.getTime() !== t1.getTime()) {
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
    animations: {
      enabled: true,
      dynamicAnimation: {
        enabled: true,
        speed: 100
      }
    },
    id: 'energy-chart',
    foreColor: '#fff',
    toolbar: {
      show: false
    },
    zoom: {
      enabled: false
    }
  },
  colors: ['var(--v-primary-base)', 'orange'],
  dataLabels: {
    enabled: false
  },
  fill: {
    colors: ['var(--v-primary-base)', 'transparent'],
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
      top: 0
      // bottom: 25
    }
  },
  legend: {
    horizontalAlign: 'right',
    itemMargin: {
      horizontal: 10 // spacing
    },
    offsetY: -5,
    onItemClick: {
      toggleDataSeries: true
    },
    onItemHover: {
      highlightDataSeries: true
    },
    position: 'top'
  },
  noData: {
    text: message.value,
    align: 'center',
    verticalAlign: 'middle',
    offsetX: 0,
    offsetY: 0,
    style: {
      color: message.value.includes('Pulling') ? 'white' : 'red',
      fontSize: '18px'
    }
  },
  stroke: {
    width: 2,
    curve: 'smooth',
    colors: ['var(--v-primary-base)', 'orange']
  },
  tooltip: {
    theme: 'dark'
  },
  xaxis: {
    labels: {
      datetimeFormatter: {
        year: 'yyyy',
        month: 'MMM \'yy',
        day: 'dd MMM',
        hour: 'HH:mm'
      }
    },
    type: 'datetime'
  },
  yaxis: {
    decimalsInFloat: 0
  }
};

const series = computed(() => {
  return [{
    name: 'Metered',
    data: meterData.value
  },
  {
    name: 'Generated',
    data: supplyData.value
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
