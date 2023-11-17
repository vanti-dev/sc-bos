<template>
  <div class="ml-3">
    <v-row class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Air Quality</h3>
      <v-spacer/>
      <v-btn class="mt-2 mr-4 elevation-0" color="neutral" @click="downloadCSV">
        Download CSV
      </v-btn>
    </v-row>

    <content-card class="mt-8 px-4 pt-6">
      <v-select
          v-model="airDevice"
          class="mb-8"
          hide-details
          :items="deviceOptions"
          item-text="label"
          item-value="value"
          label="Sensor"
          outlined/>
      <div id="air-quality-graph" :style="{width: props.width, height: props.height}">
        <div
            class="pr-5"
            :chart-options="chartOptions"
            :chart-data="chartData"
            dataset-id-key="label">
          <LineChart :chart-data="chartData" :chart-options="chartOptions"/>
        </div>
      </div>
    </content-card>
  </div>
</template>

<script setup>
import {computed, ref, watch} from 'vue';
import useAirQualityTrait from '@/composables/traits/useAirQualityTrait.js';

import LineChart from '@/components/charts/LineChart.vue';
import ContentCard from '@/components/ContentCard.vue';
import {camelToSentence} from '@/util/string';
import {DAY, HOUR, MINUTE} from '@/components/now';

const props = defineProps({
  width: {
    type: String,
    default: '100%'
  },
  height: {
    type: String,
    default: '350px'
  }
});
const airDevice = ref('');
const airQualityProps = {
  filter: () => true,
  subsystem: 'sensors',
  pollDelay: 15 * MINUTE, // 15 Minutes
  span: 24 * HOUR, // 24 Hours
  timeFrame: 30 * DAY // 30 Days
};

// Define the air quality sensor options
const {airQualitySensorHistoryValues, deviceOptions, downloadCSV} = useAirQualityTrait(airQualityProps);

// Set the default air device to the first device in the list
watch(deviceOptions, (newVal) => {
  airDevice.value = newVal[0].value;
});

// Generate the chart data
const chartData = computed(() => {
  const datasets = [];
  const labelDataMap = {};

  const airDeviceData = airQualitySensorHistoryValues[airDevice.value]?.data;
  if (!airDeviceData) {
    console.warn('No data for selected device:', airDevice.value);
    return {datasets: []};
  }

  // Transform the data to be compatible with the chart
  const remapData = (data) => {
    const transformed = [];
    for (const [key, value] of Object.entries(data.y)) {
      transformed.push({
        // Capitalize the first letter of the label
        label: camelToSentence(key).at(0).toUpperCase() + camelToSentence(key).slice(1),
        y: value,
        x: data.x
      });
    }
    return transformed;
  };

  // Generate a random color for each label
  let colorIndex = 0;
  const lightColors = [
    'rgba(255, 179, 186, 0.9)', // Light Red
    'rgba(255, 223, 186, 0.9)', // Light Orange
    'rgba(255, 255, 186, 0.9)', // Light Yellow
    'rgba(150, 255, 201, 0.9)', // Light Green
    'rgba(100, 225, 255, 1)' // Light Blue
  ];

  const randomRGBA = () => {
    const color = lightColors[colorIndex];
    colorIndex = (colorIndex + 1) % lightColors.length;
    return color;
  };


  // Aggregate data for each label across all devices
  for (const deviceData of Object.values(airQualitySensorHistoryValues)) {
    deviceData.data.forEach(dataEntry => {
      const transformedData = remapData(dataEntry);
      transformedData.forEach(({label, y, x}) => {
        if (!labelDataMap[label]) {
          labelDataMap[label] = [];
        }
        labelDataMap[label].push({y, x});
      });
    });
  }

  // Create datasets for each label
  for (const [label, data] of Object.entries(labelDataMap)) {
    const color = randomRGBA(); // Your existing randomRGBA function
    datasets.push({
      borderColor: color,
      data: data,
      fill: false,
      label: label,
      pointBackgroundColor: 'rgba(0, 0, 0, 0)',
      pointBorderColor: 'rgba(0, 0, 0, 0)',
      pointHoverBackgroundColor: 'rgb(255, 255, 255)',
      pointHoverBorderColor: color,
      pointStyle: 'circle',
      tension: 0.35
    });
  }

  return {
    datasets: datasets
  };
});

// Generate the chart options
const chartOptions = computed(() => {
  return {
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
            return date.toLocaleString('en-GB', {
              month: 'short',
              day: 'numeric',
              hour: 'numeric',
              minute: 'numeric',
              hour12: false
            });
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
          display: true,
          text: airDevice.value
        },
        min: 0
      },
      x: {
        border: {
          color: 'white'
        },
        ticks: {
          align: 'center',
          autoSkip: true,
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
        // This likely isn't needed since the data is already sorted - but it's here just in case
        // time: {
        //   displayFormats: {
        //     hour: 'HH : mm' // Format the xAxis label to 24 Hour:Minute
        //   },
        //   min: '00:00', // Set the min time to 00:00
        //   max: '23:59', // Set the max time to 23:59
        //   stepSize: 1, // Display a label for every hour
        //   unit: 'hour'
        // },
        type: 'time'
      }
    },
    type: 'line'
  };
});
</script>
