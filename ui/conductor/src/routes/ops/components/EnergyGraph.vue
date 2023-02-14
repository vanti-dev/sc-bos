<template>
  <div id="energy-graph">
    <apexchart type="area" height="180" :options="options" :series="series"/>
  </div>
</template>

<script setup>
import {computed, ref} from 'vue';

// todo: get this from the backend
const data = ref([]);
for (let i = 0; i < 24; i++) {
  data.value.push({
    x: new Date().setHours(i),
    y: 20 + Math.random() * 100
  });
}

const options = {
  chart: {
    id: 'energy-chart',
    toolbar: {show: false},
    foreColor: '#fff'
  },
  xaxis: {
    type: 'datetime'
  },
  yaxis: {
    decimalsInFloat: 0
  },
  dataLabels: {enabled: false},
  colors: ['var(--v-primary-base)'],
  stroke: {
    width: 2,
    curve: 'smooth',
    colors: ['var(--v-primary-base)']
  },
  fill: {
    colors: ['var(--v-primary-base)'],
    gradient: {
      enabled: true,
      opacityFrom: 0.55,
      opacityTo: 0.15,
      gradientToColors: ['var(--v-primary-darken1)']
    }
  },
  grid: {
    borderColor: 'var(--v-neutral-lighten2)',
    yaxis: {
      lines: {
        show: false
      }
    },
    padding: {
      top: 0,
      bottom: 25
    }
  },
  tooltip: {
    theme: 'dark'
  }
};

const series = computed(() => {
  return [{
    name: 'Energy Use',
    data: Object.values(data.value)
  }];
});

</script>

<style lang="scss" scoped>
#energy-graph {
  height: 135px;
}
</style>
