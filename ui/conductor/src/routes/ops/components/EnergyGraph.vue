<template>
  <div id="energy-graph" :style="{width, height}">
    <LineChart :chart-options="chartOptions" :chart-data="chartData" dataset-id-key="label"/>
  </div>
</template>

<script setup>
import LineChart from '@/components/charts/LineChart.vue';
import {HOUR, MINUTE, useNow} from '@/components/now';
import useMeterHistory from '@/routes/ops/components/useMeterHistory';
import useTimePeriod from '@/routes/ops/components/useTimePeriod';
import {computed} from 'vue';

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
    default: 15 * MINUTE
  },
  timeFrame: {
    type: Number,
    default: 24 * HOUR
  }
});
const {now} = useNow(() => props.span);
const {periodStart, periodEnd} = useTimePeriod(now, () => props.span, () => props.timeFrame);

const metered = useMeterHistory(() => props.metered, periodStart, periodEnd, () => props.span);
const generated = useMeterHistory(() => props.generated, periodStart, periodEnd, () => props.span);

const chartData = computed(() => {
  // Return the restructured data for the chart
  const datasets = [];
  if (props.generated) {
    datasets.push({
      borderColor: 'orange', // line color
      data: generated.seriesData.value, // data for the line
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
    });
  }

  if (props.metered) {
    datasets.push({
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
      data: metered.seriesData.value, // data for the line
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
    });
  }

  return {
    datasets
  };
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
});

/** -------------------------------------------- */
/**
 * Lifecycle hooks
 */
</script>
