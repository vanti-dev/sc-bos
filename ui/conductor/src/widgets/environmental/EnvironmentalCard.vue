<template>
  <content-card class="pt-6 pb-6">
    <v-card-title class="text-h4 mb-0">Environmental</v-card-title>
    <v-card-text class="d-flex flex-row flex-wrap justify-center align-center pa-0 white--text">
      <circular-gauge
          v-if="indoorTemperature > 0 || props.shouldWrap"
          :value="indoorTemperature"
          :color="props.gaugeColor"
          :min="tempRange.low"
          :max="tempRange.high"
          segments="30"
          class="mt-4 mx-6">
        <span class="ml-1 text-h1">
          {{ indoorTempStr }}&deg;
        </span>
        <template #title>
          Avg. Indoor Temperature
        </template>
      </circular-gauge>
      <div
          v-if="outdoorTemperature > 0 || props.shouldWrap"
          class="d-flex flex-column align-center mt-6 mx-6">
        <span class="text-h1 align-left ml-1">{{ outdoorTempStr }}&deg;</span>
        <span class="text-title text-center">External<br>Temperature</span>
      </div>
      <circular-gauge
          v-if="indoorHumidity > 0"
          :value="indoorHumidity"
          :color="props.gaugeColor"
          :min="0"
          :max="100"
          segments="30"
          class="mt-7 mx-6">
        <span class="align-baseline text-h1 ml-2">
          {{ indoorHumidityStr }}<span style="font-size: 0.7em;">%</span>
        </span>
        <template #title>
          Avg. Humidity
        </template>
      </circular-gauge>
    </v-card-text>
  </content-card>
</template>

<script setup>
import CircularGauge from '@/components/CircularGauge.vue';
import ContentCard from '@/components/ContentCard.vue';

import {useAirTemperature, usePullAirTemperature} from '@/traits/airTemperature/airTemperature.js';
import {computed} from 'vue';

const props = defineProps({
  internal: {
    type: String,
    default: null
  },
  external: {
    type: String,
    default: null
  },
  gaugeColor: {
    type: String,
    default: 'primary'
  },
  shouldWrap: {
    type: Boolean,
    default: false
  }
});

const {value: indoorValue} = usePullAirTemperature(() => props.internal);
const {
  temp: indoorTemperature,
  humidity: indoorHumidity,
  tempRange
} = useAirTemperature(indoorValue);
const {value: outdoorValue} = usePullAirTemperature(() => props.external);
const {temp: outdoorTemperature} = useAirTemperature(outdoorValue);

const vOrDash = (r) => {
  const v = r.value ?? '-';
  if (v === '-') return v;
  return v.toFixed(1);
};
const indoorTempStr = computed(() => {
  return vOrDash(indoorTemperature);
});
const indoorHumidityStr = computed(() => {
  return vOrDash(indoorHumidity);
});
const outdoorTempStr = computed(() => {
  return vOrDash(outdoorTemperature);
});

</script>
