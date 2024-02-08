<template>
  <div id="energy-graph" :style="{width, height}">
    <LineChart
        :class="props.classes"
        :chart-options="chartOptions"
        :chart-data="chartData"
        :chart-title="props.chartTitle"
        dataset-id-key="label"
        :hide-legends="props.hideLegends">
      <template #options>
        <EnergyGraphOptionsMenu
            :duration-option.sync="durationOption"
            :show-conversion.sync="showConversion"/>
      </template>
    </LineChart>
  </div>
</template>

<script setup>
import LineChart from '@/components/charts/LineChart.vue';
import {HOUR, MINUTE, useNow} from '@/components/now';
import EnergyGraphOptionsMenu from '@/routes/ops/overview/pages/widgets/energyAndDemand/EnergyGraphOptionsMenu.vue';
import useMeterHistory from '@/routes/ops/overview/pages/widgets/energyAndDemand/useMeterHistory';
import useTimePeriod from '@/routes/ops/overview/pages/widgets/energyAndDemand/useTimePeriod';
import {useCarbonIntensity} from '@/stores/carbonIntensity';
import {computed, ref} from 'vue';

const props = defineProps({
  chartTitle: {
    type: String,
    default: 'Power'
  },
  classes: {
    type: String,
    default: 'mt-n10'
  },
  color: {
    type: String,
    default: '#00bed6' // primary
  },
  colorMiddle: {
    type: String,
    default: 'rgba(51, 142, 161, 0.75)' // primary 75% opacity
  },
  hideLegends: {
    type: Boolean,
    default: false
  },
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
  }
});
const durationOption = ref({
  id: '24H',
  span: 20 * MINUTE,
  timeFrame: 24 * HOUR
});
const {now} = useNow(() => durationOption.value.span);
const {periodStart, periodEnd} = useTimePeriod(
    now,
    () => durationOption.value.span,
    () => durationOption.value.timeFrame
);
const showConversion = ref(false);
const carbonIntensity = useCarbonIntensity();
const gramsOfCO2PerKWh = ref(86);

// Simplify co2intervals computation by mapping duration IDs to carbonIntensity properties directly
const co2intervals = computed(() => {
  const mapping = {
    '24H': carbonIntensity.last24Hours,
    '1W': carbonIntensity.last7Days,
    '30D': carbonIntensity.last30Days
  };
  return mapping[durationOption.value.id] || [];
});

// Helper function to convert kWh to grams of CO2 based on the date and the carbon intensity data
const kwhToGramsOfCO2 = date => {
  const co2source = co2intervals.value;

  // Handles empty co2source array by returning the default value if no data is available
  if (!co2source.length) return gramsOfCO2PerKWh.value;

  const first = co2source[0]; // Get the first element of the array
  const last = co2source.at(-1); // Get the last element of the array

  // Handles dates before the first range
  if (date < first.from) {
    return getIntensityValue(first);
  }

  // Handles dates after the last range
  if (date > last.to) {
    return getIntensityValue(last);
  }

  // Finds and returns intensity for a date within the ranges
  const matchingRange = co2source.find(range => date >= range.from && date <= range.to);
  return matchingRange ? getIntensityValue(matchingRange) : getIntensityValue(last);
};

// Helper function to get the actual or forecasted intensity value
const getIntensityValue = (range) => range.intensity.actual ?? range.intensity.forecast;

// Set the yAxis unit based on the `showConversion` value
const yAxisUnit = computed(() => showConversion.value ? 'Grams of CO₂ / hour' : 'kW');

// Fetch the metered and generated data based on the periodStart and periodEnd values and the durationOption's span
const {seriesData: meteredSeriesData} =
    useMeterHistory(() => props.metered, periodStart, periodEnd, () => durationOption.value.span);
const {seriesData: generatedSeriesData} =
    useMeterHistory(() => props.generated, periodStart, periodEnd, () => durationOption.value.span);

// Helper function to compute the CO2 series data based on the seriesData value and the kwhToGramsOfCO2 function
const computeCO2SeriesData = seriesData => seriesData.value.map(({x, y}) => ({x, y: y * kwhToGramsOfCO2(x)}));

// Computed properties to compute the CO2 series data for the metered and generated data
const co2Metered = computed(() => computeCO2SeriesData(meteredSeriesData));
const co2Generated = computed(() => computeCO2SeriesData(generatedSeriesData));

// ----------------- Chart Data and Options ----------------- //
const chartData = computed(() => {
  const datasets = [];
  const colorEnd = 'rgba(0, 94, 107, 0.1)'; // Pre-defined for consistency and potential dynamic updates

  // Function to create a gradient; improves readability and reusability
  const createGradient = (ctx) => {
    const canvas = ctx.chart.ctx;
    const gradient = canvas.createLinearGradient(0, 0, 0, 425);
    gradient.addColorStop(0, props.color);
    gradient.addColorStop(0.5, props.colorMiddle);
    gradient.addColorStop(1, colorEnd); // Use predefined value
    return gradient;
  };

  // Helper function to avoid redundancy in dataset creation
  const addDataset = (data, isMetered = false) => {
    datasets.push({
      borderColor: isMetered ? props.color : 'orange',
      data: data,
      fill: isMetered,
      label: isMetered ? 'Metered' : 'Generated',
      backgroundColor: isMetered ? createGradient : undefined,
      pointBackgroundColor: 'rgba(0, 0, 0, 0)',
      pointBorderColor: 'rgba(0, 0, 0, 0)',
      pointHoverBackgroundColor: 'rgb(255, 255, 255)',
      pointHoverBorderColor: isMetered ? props.color : 'orange',
      pointStyle: 'circle',
      tension: 0.35
    });
  };

  // Decide which data to use based on `showConversion` and whether props are provided
  if (props.generated) {
    const generatedData = showConversion.value ? co2Generated.value : generatedSeriesData.value;
    addDataset(generatedData, false);
  }

  if (props.metered) {
    const meteredData = showConversion.value ? co2Metered.value : meteredSeriesData.value;
    addDataset(meteredData, true);
  }

  return {datasets};
});


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
          labelPointStyle: function() {
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
          text: yAxisUnit.value
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
            const date = new Date(value);

            if (date.getHours() === 0) {
              return date.toLocaleString('en-GB', {
                day: 'numeric',
                month: 'short'
              });
            } else {
              return date.toLocaleString('en-GB', {
                hour: 'numeric',
                minute: 'numeric',
                hour12: false
              });
            }
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
});
</script>
