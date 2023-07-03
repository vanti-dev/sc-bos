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
  span: {
    type: Number,
    default: 15 * 60 * 1000 // in ms
  }
});

/** -------------------------------------------- */
/**
 * ! Stage 1:
 */

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

  Object.values(seriesMap).forEach(series => {
    clearTimeout(series.handle);
  });
});

const seriesMap = reactive({
  metered: {
    baseRequest: computed(() => {
      return baseRequest(props.metered);
    }),
    data: /** @type {{x: Date, y: number}[]} */ computed(
        () => structureData(props.span, seriesMap.metered.records)
    ),
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  },
  generated: {
    baseRequest: computed(() => {
      return baseRequest(props.generated);
    }),
    data: /** @type {{x: Date, y: number}[]} */ computed(
        () => structureData(props.span, seriesMap.generated.records)
    ),
    handle: 0,
    records: /** @type {MeterReadingRecord.AsObject[]} */ []
  }
});

/**
 * Create a request for the metered series
 *
 * @param {string} name
 * @return {MeterReadingHistoryRequest.AsObject}
 */
const baseRequest = (name) => {
  if (!name) return undefined;

  // 24 hours and 15 minutes ago,
  // to account for the fact that the first reading may not be at the start of the hour
  const period = {
    startTime: new Date(now.value - (24.25 * 60 * 60 * 1000))
  };

  return {
    name,
    period,
    pageSize: 1000,
    pageToken: ''
  };
};

/**
 * Async function to poll for meter readings
 *
 * @param {*} req
 * @param {string} type
 */
async function pollReadings(req, type) {
  const all = [];
  try {
    // Loop until all meter reading records have been retrieved
    while (true) {
      // Retrieve a page of meter reading records
      const page = await listMeterReadingHistory(req, {});

      // Add the records to the 'all' array
      all.push(...page.meterReadingRecordsList);

      // Update the page token for the next page of records
      req.pageToken = page.nextPageToken;

      // If there are no more pages, break out of the loop
      if (!req.pageToken) {
        break;
      }
    }
  } catch (e) {
    console.error('error getting meter readings', e);
  }

  // Update the records and handle for the specified series type
  seriesMap[type].records = all;
  seriesMap[type].handle = setTimeout(() => pollReadings(req, type), pollDelay.value);
}

/** -------------------------------------------- */
/**
 * ! Stage 2:
 */

/**
 * Function to populate the line chart with default data
 * for the last 24 hours and an extra 15 minute in 15 minute intervals
 *
 * @param {{x: Date, y: number}[]} dataset
 * @param {number} span
 * @return {{x: Date, y: number}[]}
 */
const populateDefaultData = (dataset, span) => {
  // Populate the line chart with default data when dataset is empty
  const currentDate = new Date();

  // Adjusting current date to nearest 15 min mark
  currentDate.setMinutes(currentDate.getMinutes() - (currentDate.getMinutes() % 15), 0, 0);

  // Calculate the number of 15 min intervals in 24 hours
  const numberOfIntervals = 24.25 * 60 / 15;

  // Populate the line chart with default data
  // 24 hour divided into 15 min intervals
  for (let i = 0; i < numberOfIntervals; i++) {
    const dataPoint = {
      x: new Date(currentDate.getTime() - (i * span)),
      y: 0
    };

    dataset.unshift(dataPoint); // Update the array of objects depending on the currentDate
  }

  // Add an extra 15 min interval to the beginning of the array,
  // this is to account for the fact that the first reading may not be at the start of the hour
  dataset.unshift({
    x: new Date(dataset[0].x.getTime() - span),
    y: 0
  });

  return dataset;
};

/**
 * Function to group meter reading records by 15 minute intervals
 *
 * @param {MeterReadingRecord.AsObject[]} records
 * @return {{interval: Date, firstRecord: MeterReadingRecord.AsObject, lastRecord: MeterReadingRecord.AsObject}[]}
 */
const groupRecordsByInterval = (records) => {
  const intervalData = [];

  // Group records by 15 min intervals
  const intervalGroups = records.reduce((groups, record) => {
    const timestamp = new Date(timestampToDate(record.recordTime));
    // Round down to the nearest 15 minute interval
    const interval = new Date(Math.floor(timestamp.getTime() / props.span) * props.span);

    if (!groups[interval]) {
      groups[interval] = [];
    }

    groups[interval].push(record);

    return groups;
  }, []);

  // Convert interval keys to numerical indices
  const entries = Object.entries(intervalGroups);
  entries.forEach((entry, index) => {
    const [interval, records] = entry;


    intervalData.push({
      interval, // The 15 minute interval
      firstRecord: records[0],
      lastRecord: records[records.length - 1]
    });
  });

  return intervalData;
};


/**
 * Function to restructure the data for the line chart
 *
 * @param {number} span
 * @param {MeterReadingRecord.AsObject[]} records
 * @return {{x: Date, y: number}[]}
 */
const structureData = (span, records) => {
  let dataset = [];

  // 1.1 Populate the line chart with default data
  dataset = populateDefaultData(dataset, span);

  if (records.length) {
    // 2. Group records by 15 min intervals
    const groupedRecords = groupRecordsByInterval(records);
    // I. If there are records, update the line chart with the data
    groupedRecords.forEach(({interval, firstRecord, lastRecord}, index) => {
      const usageDiff = lastRecord.meterReading.usage - firstRecord.meterReading.usage;

      dataset.forEach((data) => {
        if (data.x == interval) {
          data.y = usageDiff;
        }
      });
    });
  }

  return dataset;
};


/** -------------------------------------------- */
/**
 * ! Stage 3:
 */


/**
 * Loop through the seriesMap and create automation with watch for:
 * pulling meter readings for the series in a specified interval (pollDelay)
 */
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

/** -------------------------------------------- */
/**
 * ! Stage 4:
 */

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
    htmlLegend: {
      // ID of the container to put the legend in
      containerID: 'legend-container'
    },
    legend: {
      display: false // Legend being overwritten by htmlLegend plugin
    },
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
    tooltip: {
      backgroundColor: '#000',
      bodyColor: '#fff',
      borderColor: '#000',
      callbacks: {
        // Format the tooltip title to Month Date, 24 Hour:Minute
        title: (data) => {
          const date = new Date(data[0].parsed.x);
          const title = date.toLocaleString('en-GB', {
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: 'numeric',
            hour12: false
          });

          return title;
        },
        labelPointStyle: function(context) {
          return {
            pointStyle: 'line',
            rotation: 0
          };
        },
        labelColor: function(context) {
          return {
            borderColor: context.dataset.borderColor,
            backgroundColor: context.dataset.borderColor,
            borderWidth: 2,
            borderDash: [2, 2]
          };
        }
      },
      borderWidth: 2,
      cornerRadius: 5,
      displayColors: true,
      enabled: true,
      interaction: {
        axis: 'xy',
        mode: 'index'
      },
      padding: 12,
      titleColor: '#fff',
      usePointStyle: true
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
          size: 12 // Specify the desired font size
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
        // autoSkip: true,
        callback: (value) => {
          // Format the xAxis label to either Month Date or 24 Hour:Minute
          let label = '';
          const date = new Date(value);

          if (date.getHours() === 0) {
            label = date.toLocaleString('en-GB', {
              day: 'numeric',
              month: 'short'
            });
          } else {
            label = date.toLocaleString('en-GB', {
              hour: 'numeric',
              minute: 'numeric',
              hour12: false
            });
          }

          return label;
        },
        color: '#fff',
        display: true,
        font: {
          size: 11 // Specify the desired font size
        },
        includeBounds: true, // Include the first and last ticks
        maxRotation: 0 // Limit xAxis label rotation to 0 degrees
      },
      grid: {
        color: ''
      },
      time: {
        displayFormats: {
          hour: 'HH : mm' // Format the xAxis label to 24 Hour:Minute
        },
        min: '00:00', // Set the min time to 00:00
        max: '23:59', // Set the max time to 23:59
        stepSize: 1, // Display a label for every hour
        unit: 'hour'
      },
      type: 'time'
    }
  },
  type: 'line'
};

// Computed property which restructures the data for the chart
const chartData = computed(() => {
  // Initialize empty arrays for labels and values
  const values = {
    metered: [],
    generated: []
  };

  // Get the data from the series prop
  const dataset = series.value;


  // Loop through each set of data in the dataset
  dataset.forEach((set, index) => {
    // Loop through each data point in the set
    set.data.forEach((data, index) => {
      // If the data point has an x value
      if (data && data.x) {
        // Create a new date object from the x value
        const newDate = new Date(data.x);
        // Get the name of the set
        const setName = set.name.toLowerCase();

        // Add the date and y value to the appropriate array
        values[setName].push({x: newDate, y: data.y});
      }
    });
  });

  // Return the restructured data for the chart
  return {
    datasets: [
      {
        borderColor: 'orange', // line color
        data: values.generated, // data for the line
        fill: false, // fill the area under the line
        label: 'Generated', // tooltip label
        mode: 'nearest', // 'index' or 'nearest
        pointBackgroundColor: 'rgba(0, 0, 0, 0)', // point background color
        pointBorderColor: 'rgba(0, 0, 0, 0)', // point border color
        pointHoverBackgroundColor: 'rgb(255, 255, 255)', // point background color on hover
        pointHoverBorderColor: 'orange', // point border color on hover
        // 'circle', 'cross', 'crossRot', 'dash', 'line', 'rect', 'rectRounded', 'rectRot', 'star', 'triangle'
        pointStyle: 'circle',
        tension: 0.35 // curve the line
      },
      {
        // Setting background gradient on metered data
        backgroundColor: (ctx) => {
          const canvas = ctx.chart.ctx;
          const gradient = canvas.createLinearGradient(0, 0, 0, 425);

          gradient.addColorStop(0, '#00bed6'); // color
          gradient.addColorStop(0.5, 'rgba(51, 142, 161, 0.75)'); // darker shade of the color
          gradient.addColorStop(1, 'rgba(0, 94, 107, 0.1)'); // almost transparent

          return gradient;
        },
        borderColor: '#00bed6', // line color
        data: values.metered, // data for the line
        fill: true, // fill area under the line graph
        label: 'Metered', // tooltip label
        mode: 'nearest', // 'index' or 'nearest
        pointBackgroundColor: 'rgba(0, 0, 0, 0)', // point background color
        pointBorderColor: 'rgba(0, 0, 0, 0)', // point border color
        pointHoverBackgroundColor: 'rgb(255, 255, 255)', // point background color on hover
        pointHoverBorderColor: '#00bed6', // point border color on hover
        // 'circle', 'cross', 'crossRot', 'dash', 'line', 'rect', 'rectRounded', 'rectRot', 'star', 'triangle'
        pointStyle: 'circle',
        tension: 0.35 // curve the line
      }
    ]
  };
});
</script>
