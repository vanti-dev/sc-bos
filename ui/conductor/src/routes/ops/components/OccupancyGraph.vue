<template>
  <div id="occupancy-graph" :style="{width, height}">
    <apexchart
        type="bar"
        height="100%"
        :options="options"
        :series="series"/>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {listOccupancySensorHistory} from '@/api/sc/traits/occupancy-history';
import {timestampToDate} from '@/api/convpb';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
  },
  width: {
    type: String,
    default: '430px'
  },
  height: {
    type: String,
    default: '275px'
  },
  span: { // how wide the bars of the histogram are / group interval
    type: Number,
    default: 15 * 60 * 1000 // in ms
  }
});

const pollDelay = computed(() => props.span / 10);
const now = ref(Date.now());
const nowHandle = ref(0);
const message = ref('No data available');
// const series = [{
//   data: [{
//     x: 'category A',
//     y: 10
//   }, {
//     x: 'category B',
//     y: 18
//   }, {
//     x: 'category C',
//     y: 13
//   },
//   {
//     x: 'category A',
//     y: 10
//   }, {
//     x: 'category B',
//     y: 18
//   }, {
//     x: 'category C',
//     y: 13
//   },
//   {
//     x: 'category A',
//     y: 10
//   }, {
//     x: 'category B',
//     y: 18
//   }, {
//     x: 'category C',
//     y: 13
//   },
//   {
//     x: 'category A',
//     y: 10
//   }, {
//     x: 'category B',
//     y: 18
//   }, {
//     x: 'category C',
//     y: 13
//   },
//   {
//     x: 'category A',
//     y: 10
//   }, {
//     x: 'category B',
//     y: 18
//   }, {
//     x: 'category C',
//     y: 13
//   }
//   ]
// }];
const seriesMap = reactive({
  [props.name]: {
    baseRequest: computed(() => {
      return baseRequest(props.name);
    }),
    data: computed(() => {
      return data(props.span, seriesMap[props.name].records);
    }),
    handle: 0,
    records: /** @type {OccupancyRecord.AsObject[]} */ []
  }
});

//
//
// Computed
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
// Methods
// Collect the data points which should be displayed on the graph
const data = (span, records) => {
  const dst = []; // Array to store data points

  if (records.length > 0) {
    // create a list of data points that show the change in value since the previous reading
    /** @type {OccupancyRecord.AsObject} */
    let lastReading = null; // Variable to store the previous reading
    /** @type {OccupancyRecord.AsObject} */
    let readingCur = null; // Variable to store the current reading

    for (const record of records) {
      if (!lastReading) {
        lastReading = record;
        readingCur = record;
        continue;
      }

      // special case if the meter was reset
      if (readingCur.occupancy.peopleCount > record.occupancy.peopleCount) {
        const diff = readingCur.occupancy.peopleCount - lastReading.occupancy.peopleCount;
        dst.push({
          x: new Date(timestampToDate(readingCur.recordTime)), // Convert timestamp to Date object
          y: diff // Store the difference in peopleCount
        });
        lastReading = readingCur = record;
        continue;
      }

      readingCur = record;
      const t0 = timestampToDate(lastReading.recordTime); // Convert timestamp to Date object
      const t1 = timestampToDate(record.recordTime); // Convert timestamp to Date object
      const d = t1 - t0; // Calculate time difference in milliseconds

      // Check time difference against interval
      if (d > span) {
        // Calculate the number of segments within the time difference
        const segmentCount = Math.floor(d / span);
        // Calculate the difference per segment
        const diff = (record.occupancy.peopleCount - lastReading.occupancy.peopleCount) / segmentCount;
        lastReading = record;
        dst.push({
          x: new Date(t1),
          y: diff
        });
      }
    }

    // process the last reading, if we haven't already
    const finalReading = records[records.length - 1];
    const t0 = timestampToDate(lastReading.recordTime); // Convert timestamp to Date object
    const t1 = timestampToDate(finalReading.recordTime);
    if (t0.getTime() !== t1.getTime()) {
      const diff = finalReading.occupancy.peopleCount - lastReading.occupancy.peopleCount;
      dst.push({
        x: new Date(t1),
        y: diff
      });
    }
  }

  // Reduce data points by interval and update with maximum peopleCount
  const reducedDst = [];
  const interval = span; // 15 minutes in milliseconds

  for (const dataPoint of dst) {
    const timestamp = dataPoint.x.getTime(); // Get the timestamp in milliseconds
    const hourStart = Math.floor(timestamp / 3600000) * 3600000; // Get the start of the hour in milliseconds
    const intervalStart = Math.floor(
        (timestamp - hourStart) / interval) * interval; // Get the start of the interval in milliseconds

    const newTimestamp = hourStart + intervalStart; // Update the timestamp to the start of the interval

    // Find existing data point with the same timestamp
    const existingDataPoint = reducedDst.find(dp => dp.x.getTime() === newTimestamp);

    if (existingDataPoint) {
      // Take the maximum peopleCount value within the interval
      existingDataPoint.y = Math.max(existingDataPoint.y, dataPoint.y);
    } else {
      reducedDst.push({
        x: new Date(newTimestamp),
        y: dataPoint.y
      });
    }
  }

  return reducedDst;
};


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

/**
 *
 * @param {*} req
 * @param {string} type
 */
async function pollReadings(req, type) {
  const all = [];
  try {
    while (true) {
      const page = await listOccupancySensorHistory(req, {});

      all.push(...page.occupancyRecordsList);
      req.pageToken = page.nextPageToken;
      if (!req.pageToken) {
        break;
      }
      message.value = 'No data available';
    }
  } catch (e) {
    console.error('error getting occupancy readings', e);
  }

  seriesMap[type].records = all;
  seriesMap[type].handle = setTimeout(() => pollReadings(req, type), pollDelay.value);
}


//
//
// Bar graph styling and other options
const options = {
  chart: {
    animations: {
      enabled: true,
      dynamicAnimation: {
        enabled: true,
        speed: 100
      }
    },
    id: 'occupancy-chart',
    foreColor: '#fff',
    toolbar: {
      show: false
    }
  },
  dataLabels: {
    enabled: false
  },
  fill: {
    colors: ['#C5CC3CBF']
  },
  grid: {
    borderColor: 'var(--v-neutral-lighten2)',
    yaxis: {
      lines: {
        show: false
      }
    },
    xaxis: {
      lines: {
        show: false
      }
    },
    padding: {
      top: 0
    }
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
  plotOptions: {
    bar: {
      horizontal: false,
      startingShape: 'flat',
      endingShape: 'rounded'
    }
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


//
//
// Lifecycle
onMounted(() => {
  nowHandle.value = setInterval(() => {
    now.value = Date.now();
  }, pollDelay.value);
});

onUnmounted(() => {
  clearInterval(nowHandle.value);
});
</script>

<style lang="scss" scoped>
</style>
