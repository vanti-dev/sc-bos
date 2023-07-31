<template>
  <div id="occupancy-graph" :style="{ width, height }">
    <BarChart :chart-data="chartData" :chart-options="chartOptions"/>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import BarChart from '@/components/charts/BarChart.vue';
import useOccupancyData from '@/routes/ops/components/useOccupancyData.js';
import {generateAreaBackground} from '@/components/charts/chartHelpers';

const props = defineProps({
  name: {
    type: String,
    default: 'building'
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
    // how wide the bars of the histogram are / group interval
    type: Number,
    default: 15 * 60 * 1000 // in ms
  }
});

const {chartSeries, chartSegments} = useOccupancyData(props);

const setAnnotation = computed(() => {
  const segments = chartSegments.value;
  const conditionValue = 'occupied';
  const mainColor = 'rgba(125, 125, 125, 0.25)';
  const secondaryColor = 'transparent';

  return generateAreaBackground(segments, conditionValue, mainColor, secondaryColor);
});

//
//
// Bar graph styling and other options
const chartOptions = computed(() => {
  return {
    animation: {
      duration: 500
    },
    layout: {
      padding: {
        top: 0,
        right: 10,
        bottom: 20,
        left: 10
      }
    },
    maintainAspectRatio: false,
    mode: 'none',
    responsive: true,
    plugins: {
      annotation: setAnnotation.value,
      legend: {
        display: false
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
        min: 0,
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
        }
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
            size: 10 // Specify the desired font size
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
});

const chartData = computed(() => {
  const dataset = chartSeries.value[0].data;

  let labels = [];

  labels = dataset.map((data) => {
    const newDate = new Date(data.x);
    // removing seconds from the time
    const hour = newDate.getHours().toString().padStart(2, '0');
    const minute = newDate.getMinutes().toString().padStart(2, '0');
    const formattedTime = hour + ':' + minute;

    return formattedTime;
  });

  const data = dataset.map((data) => data.y);

  // returning restructured data for the bar chart
  return {
    labels, // collection of time intervals
    datasets: [
      {
        type: 'bar', // bar chart type
        label: 'Occupancy', // tooltip label
        maxBarThickness: 20,
        data,
        backgroundColor: '#C5CC3CBF',
        borderColor: '#fff',
        borderRadius: 5
      }
    ]
  };
});
</script>
