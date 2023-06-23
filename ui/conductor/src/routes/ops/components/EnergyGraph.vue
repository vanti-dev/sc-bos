<template>
  <div id="energy-graph" :style="{width, height}">
    <LineChart :chart-options="chartOptions" :chart-data="chartData"/>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {listMeterReadingHistory} from '@/api/sc/traits/meter-history';
import LineChart from '@/components/charts/LineChart.vue';
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
    default: '100%'
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

const seriesMap = reactive({
  metered: {
    baseRequest: computed(() => {
      return baseRequest(props.metered);
    }),
    data: [],
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  },
  generated: {
    baseRequest: computed(() => {
      return baseRequest(props.generated);
    }),
    data: [],
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  }
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


// Function to calculate the difference between two data points
const calculateDifference = (lastReading, record) => {
  return record.meterReading.usage - lastReading.meterReading.usage;
};

// Function to add a data point for the difference if the meter was reset
const addDataPointForReset = (lastReading, records, dataPoints) => {
  const diff = calculateDifference(records[0], lastReading);
  dataPoints.push({x: new Date(timestampToDate(lastReading.recordTime)), y: diff});
};

// Function to add data points for each segment if the time difference is greater than the specified span
const addDataPointsForSegment = (lastReading, record, span, dataPoints) => {
  const segmentCount = Math.floor(
      (timestampToDate(record.recordTime) - timestampToDate(lastReading.recordTime)) / span
  );
  const diff = calculateDifference(lastReading, record) / segmentCount;

  for (let i = 1; i <= segmentCount; i++) {
    const x = new Date(timestampToDate(lastReading.recordTime).getTime() + i * span);
    const y = diff;
    dataPoints.push({x, y});
  }
};

// Function to add a data point for the final reading if it hasn't already been added
const addDataPointForFinalReading = (lastReading, records, dataPoints) => {
  const finalReading = records[records.length - 1];
  const [t0, t1] = [timestampToDate(lastReading.recordTime), timestampToDate(finalReading.recordTime)];

  if (t0.getTime() !== t1.getTime()) {
    const diff = calculateDifference(lastReading, finalReading);
    dataPoints.push({x: new Date(t1), y: diff});
  }
};

// Function to collect the data points which should be displayed on the graph
const data = (span, records) => {
  // If there are no records, return an empty array
  if (records.length === 0) {
    return [];
  }

  // Initialize an empty array to hold the data points
  const dataPoints = [];

  // Initialize the last reading to the first record
  let lastReading = records[0];

  // Iterate through the records and calculate the data points
  for (const record of records.slice(1)) {
    // If the meter was reset, add a data point for the difference
    if (lastReading.meterReading.usage > record.meterReading.usage) {
      addDataPointForReset(lastReading, records, dataPoints);
      lastReading = record;
      continue;
    }

    // Calculate the time difference between the last and current records
    const d = timestampToDate(record.recordTime) - timestampToDate(lastReading.recordTime);

    // If the time difference is greater than the specified span, add data points for each segment
    if (d > span) {
      addDataPointsForSegment(lastReading, record, span, dataPoints);
      lastReading = record;
    }
  }

  // Add a data point for the final reading if it hasn't already been added
  addDataPointForFinalReading(lastReading, records, dataPoints);

  // Return the data points array
  return dataPoints;
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

  seriesMap[type].records = all;
  seriesMap[type].handle = setTimeout(() => pollReadings(req, type), pollDelay.value);
}


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
// Chart Options and Data
const chartOptions = {
  animation: {
    duration: 500
  },
  hover: {
    intersect: false
  },
  layout: {
    padding: {
      top: 0,
      right: 10,
      bottom: 0,
      left: 10
    }
  },
  maintainAspectRatio: false,
  responsive: true,
  plugins: {
    title: {
      display: true,
      text: '',
      color: '#f0f0f0',
      font: {
        size: 16
      },
      padding: {
        bottom: 0
      }
    },
    legend: {
      display: true
    },
    tooltip: {
      backgroundColor: 'rgba(0, 0, 0, 1)',
      padding: 12,
      cornerRadius: 5,
      borderColor: '#000',
      borderWidth: 2,
      titleColor: '#fff',
      bodyColor: '#fff',
      displayColors: false
    }
  },
  scales: {
    y: {
      border: {
        color: 'white'
      },
      grid: {
        color: 'rgba(100, 100, 100, 0.35)'
      },
      ticks: {
        color: '#fff',
        display: true,
        font: {
          size: 12// Specify the desired font size
        }
      },
      title: {
        display: false,
        text: ''
      },
      min: 0
    },
    x: {
      border: {
        color: 'white'
      },
      ticks: {
        align: 'center',
        color: '#fff',
        display: true,
        font: {
          size: 10// Specify the desired font size
        },
        // Limit xAxis label rotation to 0 degrees
        maxRotation: 0,
        minRotation: 0
      },
      grid: {
        color: ''
      }
    }
  }
};


// Define a computed property that restructures the data for the chart
const chartData = computed(() => {
  // Initialize empty arrays for labels and values
  let labels = {
    metered: [],
    generated: []
  };

  const values = {
    metered: [],
    generated: []
  };

  // Get the data from the series prop
  const dataset = series.value;

  // Define a function to format the time
  const timeFormat = (time) => {
    const hour = time.getHours().toString().padStart(2, '0');
    const minute = time.getMinutes().toString().padStart(2, '0');
    const formattedTime = hour + ':' + minute;

    return formattedTime;
  };

  // Loop through each set of data in the dataset
  dataset.forEach((set, index) => {
    // Loop through each data point in the set
    set.data.forEach((data, index) => {
      // If the data point has an x value
      if (data && data.x) {
        // Create a new date object from the x value
        const newDate = new Date(data.x);
        // Format the time using the timeFormat function
        const formattedTime = timeFormat(newDate);

        // If the set is for metered data, add the time and value to the metered arrays
        if (set.name === 'Metered') {
          labels.metered.push(formattedTime);
          values.metered.push(data.y);
        } else {
          // Otherwise, add the time and value to the generated arrays
          labels.generated.push(formattedTime);
          values.generated.push(data.y);
        }
      }
    });
  });

  // Reduce duplicate values in the labels array
  labels = [...labels.metered, ...labels?.generated].reduce((acc, curr) => {
    if (acc.indexOf(curr) === -1) {
      acc.push(curr);
    }
    return acc;
  }, []);

  // Return the restructured data for the chart
  return {
    labels, // collection of time intervals
    datasets: [
      {
        borderColor: 'orange',
        data: values.generated,
        fill: false,
        label: 'Generated', // tooltip label
        mode: 'nearest', // 'index' or 'nearest
        pointBackgroundColor: 'rgba(0, 0, 0, 0)',
        pointBorderColor: 'rgba(0, 0, 0, 0)',
        pointHoverBackgroundColor: 'rgb(255, 255, 255)',
        pointHoverBorderColor: 'orange',
        pointStyle: 'circle',
        tension: 0.35
      },
      {
        // Setting background gradient on metered data
        backgroundColor: (ctx) => {
          const canvas = ctx.chart.ctx;
          const gradient = canvas.createLinearGradient(0, 0, 0, 425);

          gradient.addColorStop(0, '#00bed6');
          gradient.addColorStop(0.5, 'rgba(51, 142, 161, 0.75)');
          gradient.addColorStop(1, 'rgba(0, 94, 107, 0.1)'); // Darker shade of the color

          return gradient;
        },
        borderColor: '#00bed6',
        data: values.metered, // actual data
        fill: true,
        label: 'Metered', // tooltip label
        mode: 'nearest', // 'index' or 'nearest
        pointBackgroundColor: 'rgba(0, 0, 0, 0)',
        pointBorderColor: 'rgba(0, 0, 0, 0)',
        pointHoverBackgroundColor: 'rgb(255, 255, 255)',
        pointHoverBorderColor: '#00bed6',
        pointStyle: 'circle',
        tension: 0.35
      }
    ]
  };
});


//
//
// Watcher
Object.entries(seriesMap).forEach(([name, series]) => {
  // Watching baseRequest for changes
  watch(() => series.baseRequest, (request) => {
    clearTimeout(series.handle);
    series.records = [];

    // create new stream
    if (request) {
      pollReadings(request, name);
    }
  }, {immediate: true, deep: true, flush: 'sync'});
});

Object.entries(seriesMap).forEach(([name, series]) => {
  // Watching records for changes
  watch(() => series.records, (records) => {
    const seriesCollection = data(props.span, records);

    if (series.data.length === 0) {
      series.data = seriesCollection;
    } else {
      series.data.push(seriesCollection[seriesCollection.length - 1]);
      series.data.shift();
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

  Object.values(seriesMap).forEach(series => {
    clearTimeout(series.handle);
  });
});
</script>

<style lang="scss" scoped>
</style>
