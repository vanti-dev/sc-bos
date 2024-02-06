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
          label="Data source"
          :loading="isFetching"
          outlined/>
      <div class="pt-4" id="air-quality-graph" :style="{width: props.width, height: props.height}">
        <div
            class="pr-5"
            :chart-options="chartOptions"
            :chart-data="chartData"
            dataset-id-key="label">
          <LineChart :chart-data="chartData" :chart-options="chartOptions"/>
        </div>
      </div>

      <!-- Most recent values -->
      <v-row class="d-flex flex-row justify-space-between mt-5 mb-4 ml-15 mr-4">
        <v-col
            v-for="(recentValue, key) in filteredValues"
            :key="key"
            cols="auto"
            class="text-h1 align-self-auto"
            style="line-height: 0.35em;">
          {{ recentValue.value }}
          <span style="font-size: 0.5em;">{{ recentValue.unit }}</span><br>
          <span
              class="text-h6"
              :style="{lineHeight: '0.4em', color: getColorForKey(key)}">
            {{ recentValue.label }}
          </span>
        </v-col>
      </v-row>
    </content-card>
  </div>
</template>

<script setup>
import LineChart from '@/components/charts/LineChart.vue';
import ContentCard from '@/components/ContentCard.vue';
import {HOUR, MINUTE} from '@/components/now';
import useAirQualityTrait from '@/composables/traits/useAirQualityTrait.js';
import {camelToSentence} from '@/util/string';
import {computed, reactive} from 'vue';

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
const airQualityProps = reactive({
  filter: () => true, // Filter function to filter out data in deviceData
  name: '', // Zone or device name - if `subsystem` is not set, the name must be set
  pollDelay: 15 * MINUTE, // 15 Minutes
  span: 15 * MINUTE, // 24 Hours
  subsystem: 'zones', // If `name` is not set, the subsystem must be set
  timeFrame: 24 * HOUR // 24 Hours
});

// Define the air quality trait options
const {
  acronyms,
  airQualitySensorHistoryValues,
  airDevice,
  deviceOptions,
  downloadCSV,
  isFetching,
  readComfortValue
} = useAirQualityTrait(airQualityProps);

// Define a color mapping for each key
const colorMapping = {
  carbonDioxideLevel: 'rgba(255, 179, 186, 0.9)', // Light Red
  volatileOrganicCompounds: 'rgba(255, 223, 186, 0.9)', // Light Orange
  airPressure: 'rgba(255, 255, 186, 0.9)', // Light Yellow
  comfort: 'rgba(150, 255, 201, 0.9)', // Light Green
  infectionRisk: 'rgba(100, 225, 255, 1)', // Light Blue,
  score: 'rgba(100, 100, 255, 1)' // Purple-ish
};

// Function to get color based on the key
const getColorForKey = (key) => {
  return colorMapping[key] || 'rgba(0, 0, 0, 0.5)'; // Default color
};

// Returns the most recent (last) values for each key
const mostRecentValues = computed(() => {
  const airDeviceData = airQualitySensorHistoryValues[airDevice.value]?.data;
  const mostRecentValues = {};

  if (airDeviceData) {
    const mostRecent = airDeviceData[airDeviceData.length - 1];
    if (!mostRecent) {
      return {};
    }

    for (const [key, value] of Object.entries(mostRecent?.y)) {
      mostRecentValues[key] = value;
    }
  }

  return mostRecentValues;
});

// Filtering out values that are null
const filteredValues = computed(() => {
  return Object.entries(mostRecentValues.value).reduce((acc, [key, value]) => {
    const showValue = showHideValue(value, key);
    if (showValue.value) {
      acc[key] = showValue;
    }
    return acc;
  }, {});
});

// Function to show/hide values based on the key
const showHideValue = (value, key) => {
  const unit = acronyms[key].unit;
  const label = acronyms[key].label;

  if (key === 'comfort') {
    return {
      label,
      unit,
      value: readComfortValue(value) !== 'Unspecified' ? readComfortValue(value) : null
    };
  } else {
    // For other keys, check if the value is greater than 0
    return {
      label,
      unit,
      value: value > 0 ? value.toFixed(2) : null
    };
  }
};


// Generate the chart data
const chartData = computed(() => {
  const datasets = [];
  const labelDataMap = {};

  const airDeviceData = airQualitySensorHistoryValues[airDevice.value]?.data;

  if (airDeviceData) {
    // Transform the data to be compatible with the chart
    const remapData = (data) => {
      const transformed = [];
      for (const [key, value] of Object.entries(data.y)) {
        transformed.push({
          // Capitalize the first letter of the label
          label: key,
          y: value,
          x: data.x
        });
      }
      return transformed;
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
    for (const [labelKey, data] of Object.entries(labelDataMap)) {
      const color = getColorForKey(labelKey); // Use the color from the mapping

      // Check if the labelKey has a corresponding label in acronyms; otherwise, use camelToSentence
      const label = acronyms[labelKey]?.label || camelToSentence(labelKey);
      const formattedLabel = label.charAt(0).toUpperCase() + label.slice(1);

      datasets.push({
        borderColor: color,
        data: data,
        fill: false,
        label: formattedLabel,
        pointBackgroundColor: 'rgba(0, 0, 0, 0)',
        pointBorderColor: 'rgba(0, 0, 0, 0)',
        pointHoverBackgroundColor: 'rgb(255, 255, 255)',
        pointHoverBorderColor: color,
        pointStyle: 'circle',
        tension: 0.35
      });
    }
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
          },
          // Format the label - if it is comfort, format the value; otherwise, format the label and value
          label: (context) => {
            const label = context.dataset.label;
            const value = context.parsed.y;

            // Find the matching acronym.label
            const acronymLabel = Object.entries(acronyms).find(([key, value]) => {
              return value.label === label;
            })?.[1]?.label;

            const acronymUnit = Object.entries(acronyms).find(([key, value]) => {
              return value.label === label;
            })?.[1]?.unit;

            // If the label is comfort, format the value
            if (label === 'Comfort') {
              return `${acronymLabel}: ${readComfortValue(value)}`;
            } else {
              return `${acronymLabel}: ${value.toFixed(2)} ${acronymUnit}`;
            }
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
