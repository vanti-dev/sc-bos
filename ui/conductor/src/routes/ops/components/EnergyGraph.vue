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
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';

const props = defineProps({
  metered: {
    type: String,
    default: 'building'
  },
  generated: {
    type: String,
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
const baseRequest = (name) => {
  if (!name) return undefined;

  const period = {
    startTime: new Date(now.value - 24 * 60 * 60 * 1000)
  };

  return {
    name,
    period,
    pageSize: 1000,
    pageToken: ''
  };
};

// Collect the data points which should be displayed on the graph
const data = (span, records) => {
  const dst = [];

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
    if (t0.getTime() !== t1.getTime()) {
      const diff = finalReading.meterReading.usage - lastReading.meterReading.usage;
      dst.push({
        x: new Date(t1),
        y: diff
      });
    }
  }
  return dst;
};


const seriesMap = reactive({
  metered: {
    baseRequest: computed(() => {
      return baseRequest(props.metered);
    }),
    data: computed(() => {
      return data(props.span, seriesMap.metered.records);
    }),
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  },
  generated: {
    baseRequest: computed(() => {
      return baseRequest(props.generated);
    }),
    data: computed(() => {
      return data(props.span, seriesMap.generated.records);
    }),
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  }
});

// Graph status message
const message = ref('');

/**
 *
 * @param {*} req
 * @param {string} type
 */
async function pollReadings(req, type) {
  const all = [];
  try {
    while (true) {
      const page = await listMeterReadingHistory(req, {});

      all.push(...page.meterReadingRecordsList);
      req.pageToken = page.nextPageToken;
      if (!req.pageToken) {
        break;
      }
      message.value = 'No data available';
    }
  } catch (e) {
    console.error('error getting meter readings', e);
  }

  seriesMap[type].records = all;
  seriesMap[type].handle = setTimeout(() => pollReadings(req, type), pollDelay.value);
}

onUnmounted(() => {
  Object.values(seriesMap).forEach(series => {
    clearTimeout(series.handle);
  });
});

//
//
// Computed
// Generating object with name and data key/values required for graph 'series' prop
const series = computed(() => {
  return Object.entries(seriesMap).map(([seriesName, seriesData]) => {
    const capitalisedName = seriesName.charAt(0).toUpperCase() + seriesName.slice(1);
    const data = seriesMap[seriesName].data;

    if (data && data.length > 0) {
      return {name: capitalisedName, data};
    } else {
      return null;
    }
  }).filter(obj => obj !== null);
});


//
//
// Line graph styling and other options
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


//
//
// Watcher
Object.entries(seriesMap).forEach(([name, series]) => {
  watch(() => series.baseRequest, (request) => {
    message.value = 'Pulling data';
    clearTimeout(series.handle);
    series.records = [];

    // create new stream
    if (request) {
      pollReadings(request, name);
    }
  }, {immediate: true, deep: true, flush: 'sync'});
});
</script>

<style lang="scss" scoped>
#energy-graph {
  height: 135px;
}
</style>
