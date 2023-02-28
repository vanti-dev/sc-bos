<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Environmental</h4>
    <div class="d-flex flex-column flex-md-row flex-lg-column">
      <circular-gauge
          :value="temperature"
          min="15"
          max="35"
          segments="30"
          class="align-self-center mb-6 mb-md-0 mb-lg-8 mr-md-8 mr-lg-0">
        {{ temperature.toFixed(1) }}&deg;
        <template #title>Avg. Indoor Temperature</template>
      </circular-gauge>
      <div class="align-self-center mb-6 mb-md-0 mb-lg-8 mr-md-8 mr-lg-0" style="width: 205px;">
        <span
            class="text-title"
            style="display: inline-block; width: 100px;">External Temperature</span>
        <span
            class="text-h1"
            style="display: inline-block; float: right;">{{ externalTemperature.toFixed(1) }}&deg;</span>
      </div>
      <circular-gauge
          :value="humidity"
          max="100"
          segments="30"
          class="align-self-center">
        <span class="align-baseline">
          {{ humidity.toFixed(1) }}<span style="font-size: 0.7em;">&percnt;</span>
        </span>
        <template #title>Avg. Humidity</template>
      </circular-gauge>
    </div>
  </content-card>
</template>

<script setup>

import ContentCard from '@/components/ContentCard.vue';
import CircularGauge from '@/components/CircularGauge.vue';

const props = defineProps({
  temperature: {
    type: Number,
    default: 0
  },
  humidity: {
    type: Number,
    default: 0
  },
  externalTemperature: {
    type: Number,
    default: 0
  }
});

</script>

<style scoped>

</style>
