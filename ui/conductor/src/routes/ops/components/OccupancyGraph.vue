<template>
  <apexchart
      type="bar"
      height="200px"
      :options="options"
      :series="series"/>
</template>

<script setup>
import {ref} from 'vue';

const message = ref('No data available');


const series = [{
  data: [{
    x: 'category A',
    y: 10
  }, {
    x: 'category B',
    y: 18
  }, {
    x: 'category C',
    y: 13
  },
  {
    x: 'category A',
    y: 10
  }, {
    x: 'category B',
    y: 18
  }, {
    x: 'category C',
    y: 13
  },
  {
    x: 'category A',
    y: 10
  }, {
    x: 'category B',
    y: 18
  }, {
    x: 'category C',
    y: 13
  },
  {
    x: 'category A',
    y: 10
  }, {
    x: 'category B',
    y: 18
  }, {
    x: 'category C',
    y: 13
  },
  {
    x: 'category A',
    y: 10
  }, {
    x: 'category B',
    y: 18
  }, {
    x: 'category C',
    y: 13
  }
  ]
}];

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
  // colors: ['var(--v-primary-base)', 'orange'],
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
      // bottom: 25
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
      horizontal: false
    }
  },
  // stroke: {
  //   width: 2,
  //   curve: 'smooth',
  //   colors: ['var(--v-primary-base)', 'orange']
  // },
  tooltip: {
    theme: 'dark'
  },
  // xaxis: {
  //   labels: {
  //     datetimeFormatter: {
  //       year: 'yyyy',
  //       month: 'MMM \'yy',
  //       day: 'dd MMM',
  //       hour: 'HH:mm'
  //     }
  //   },
  //   type: 'datetime'
  // },
  yaxis: {
    decimalsInFloat: 0
  }
};
</script>
