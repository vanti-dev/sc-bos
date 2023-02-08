<template>
  <v-container fluid class="pa-0 pt-2 d-flex flex-column" style="min-width: 270px; max-width: 1200px">
    <div class="d-flex flex-column flex-md-row">
      <h3 class="text-h3 pt-2 pb-6 flex-grow-1">Status: Building Overview</h3>
      <content-card class="mr-6 mb-6 mr-sm-0" style="min-width: 270px;">
        <div class="py-1" style="text-align: center">
          <span class="text-title">Smart Core OS: </span>
          <span class="text-title-bold text-uppercase success--text text--lighten-4">Online</span>
        </div>
      </content-card>
    </div>
    <div class="d-flex flex-column flex-lg-row">
      <div class="flex-grow-1 d-flex flex-column mr-lg-6">
        <content-card class="mb-6 mr-0 mr-lg-6 px-6">
          <h4 class="text-h4 order-md-last order-sm-first pb-4 pb-md-0 pt-0 pt-md-4">System Monitor</h4>
        </content-card>
        <content-card class="mb-6 mr-0 mr-lg-6 px-6">
          <h4 class="text-h4 order-md-last order-sm-first pb-4 pb-md-0 pt-0 pt-md-4">Energy</h4>
        </content-card>
      </div>
      <div class="d-flex flex-column" style="min-width: 270px;">
        <content-card class="mb-6 mr-0 px-6">
          <h4 class="text-h4 order-md-last order-sm-first pb-4 pb-md-0 pt-0 pt-md-4">Occupancy</h4>
        </content-card>
        <content-card class="mb-6 d-flex flex-column px-6 pt-md-6">
          <h4 class="text-h4 order-md-last order-sm-first pb-4 pb-md-0 pt-0 pt-md-4">Environmental</h4>
          <div class="d-flex flex-column flex-md-row flex-lg-column">
            <circular-gauge
                :value="temperature"
                min="-10"
                max="30"
                segments="30"
                class="align-self-center mb-8 mb-md-0 mb-lg-8 mr-md-8 mr-lg-0">
              {{ temperature.toFixed(1) }}&deg;
              <template #title>Avg. Indoor Temperature</template>
            </circular-gauge>
            <circular-gauge
                :value="humidity"
                max="100"
                segments="30"
                class="align-self-center">
              {{ humidity.toFixed(1) }}&percnt;
              <template #title>Avg. Humidity</template>
            </circular-gauge>
          </div>
          <v-slider v-model="sliderVal" min="-10" max="30" step="0.5"/>
        </content-card>
      </div>
    </div>
  </v-container>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import CircularGauge from '@/components/CircularGauge.vue';
import {computed, ref} from 'vue';

const sliderVal = ref(0);
const temperature = computed(() => sliderVal.value);
const humidity = computed(() => ((sliderVal.value+10)*100/40));

const floors = [
  {
    name: 'G',
    title: 'Ground Floor',
    zone: 'L00'
  },
  {
    name: '1',
    title: '1st Floor',
    zone: 'L01'
  },
  {
    name: '2',
    title: '2nd Floor',
    zone: 'L02'
  },
  {
    name: '3',
    title: '3rd Floor',
    zone: 'L03'
  },
  {
    name: '4',
    title: '4th Floor',
    zone: 'L04'
  },
  {
    name: '5',
    title: '5th Floor',
    zone: 'L05'
  },
  {
    name: '6',
    title: '6th Floor',
    zone: 'L06'
  },
  {
    name: '7',
    title: '7th Floor',
    zone: 'L07'
  },
  {
    name: '8',
    title: '8th Floor',
    zone: 'L08'
  },
  {
    name: '9',
    title: '9th Floor',
    zone: 'L09'
  },
  {
    name: '10',
    title: '10th Floor',
    zone: 'L10'
  }
];

</script>
