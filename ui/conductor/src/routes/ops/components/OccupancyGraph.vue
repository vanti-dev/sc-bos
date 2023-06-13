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

const seriesMap = reactive({
  occupancy: {
    baseRequest: computed(() => {
      return baseRequest(props.name);
    }),
    data: computed(() => {
      return data(seriesMap.occupancy.records);
    }),
    handle: 0,
    records: /** @type {OccupancyRecord.AsObject[]} */ []
  }
});

//
//
// Computed
// Formatting the data for the bar chart
const series = computed(() => {
  return Object.entries(seriesMap).map(([seriesName, seriesData]) => {
    const data = seriesMap[seriesName].data;

    if (data && data.length > 0) {
      return {name: 'Occupancy', data};
    } else {
      return null;
    }
  }).filter(obj => obj !== null);
});


//
//
// Methods
const data = (records) => {
  const intervalsMap = [];
  const currentDate = new Date();

  // Adjusting current date to nearest half-hour mark
  currentDate.setMinutes(currentDate.getMinutes() - (currentDate.getMinutes() % 30), 0, 0);

  // Populate the bar chart with data
  // 24 hour divided into 30 min intervals
  for (let i = 0; i < 48; i++) {
    const dataPoint = {
      x: new Date(currentDate.getTime() - (i * 30 * 60 * 1000)),
      y: 0
    };

    intervalsMap.unshift(dataPoint); // Update the array of objects depending on the currentDate
  }

  // Split each hour into 30min intervals and group the records while finding the highest value
  for (const record of records) {
    const recordTime = new Date(timestampToDate(record.recordTime));
    const minute = recordTime.getMinutes();
    const intervalStart = new Date(recordTime); // Separating the hours into 30 min intervals
    intervalStart.setMinutes(minute < 30 ? 0 : 30, 0, 0); // Start of the interval

    // Looking for the existing interval record
    const existingInterval = intervalsMap.find(
        intrvl => intrvl.x.getTime() === intervalStart.getTime()
    );

    // Updating the interval record if higher peopleCount comes in
    if (record.occupancy.peopleCount > existingInterval.y) {
      existingInterval.y = record.occupancy.peopleCount;
    }
  }

  return intervalsMap;
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
    },
    zoom: {
      enabled: false
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
        show: true
      }
    },
    padding: {
      top: 0,
      bottom: 25
    }
  },
  plotOptions: {
    bar: {
      horizontal: false
    }
  },
  tooltip: {
    theme: 'dark',
    x: {
      format: 'dd MMM yyyy',
      formatter: function(value, {series, seriesIndex, dataPointIndex, w}) {
        const newDate = new Date(value);
        return newDate.toLocaleString();
      }
    }
  },
  xaxis: {
    labels: {
      datetimeFormatter: {
        year: 'yyyy',
        month: 'MMM \'yy',
        day: 'dd MMM',
        hour: 'HH:mm'
      },
      datetimeUTC: false,
      hideOverlappingLabels: true,
      showDuplicates: false,
      trim: false
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
