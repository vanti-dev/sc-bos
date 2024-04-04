<template>
  <content-card class="mb-5 d-flex flex-column pt-7 pb-0">
    <h4 class="text-h4 pl-4 pb-8 pt-1">Environmental</h4>
    <div class="d-flex flex-column align-center mb-4">
      <v-col cols="auto" class="ma-0 pa-0">
        <circular-gauge
            v-if="indoorTemperature > 0 || props.shouldWrap"
            :value="indoorTemperature"
            :color="props.gaugeColor"
            :min="tempRange.low"
            :max="tempRange.high"
            segments="30"
            style="max-width: 140px;"
            class="mt-2 mb-5 ml-3 mr-2">
          <span class="mt-n4 ml-1 text-h1">
            {{ indoorTempStr }}&deg;
          </span>
          <template #title>
            <span class="ml-n1 mb-2">Avg. Indoor Temperature</span>
          </template>
        </circular-gauge>
      </v-col>
      <v-col cols="auto" class="mt-auto mb-0 pb-2 px-0">
        <div
            v-if="outdoorTemperature > 0 ||
              props.shouldWrap
            "
            :class="[indoorHumidity > 0 ? 'mb-7' : 'mb-2',
                     'd-flex flex-column align-center ml-2']"
            style="width: 150px;">
          <span
              class="text-h1 align-left mb-3"
              style="display: inline-block;">{{ outdoorTempStr }}&deg;
          </span>
          <span
              class="text-title text-center"
              style="display: inline-block; width: 100px;">
            External Temperature
          </span>
        </div>
      </v-col>
      <v-col cols="auto" class="pa-0">
        <circular-gauge
            v-if="indoorHumidity > 0"
            :value="indoorHumidity"
            :color="props.gaugeColor"
            :min="0"
            :max="100"
            segments="30"
            style="max-width: 140px;"
            class="mt-2">
          <span class="align-baseline text-h1 mt-n2">
            {{ indoorHumidityStr }}<span style="font-size: 0.7em;">%</span>
          </span>
          <template #title>
            <span class="mb-2">Avg. Humidity</span>
          </template>
        </circular-gauge>
      </v-col>
    </div>
  </content-card>
</template>

<script setup>
import CircularGauge from '@/components/CircularGauge.vue';
import ContentCard from '@/components/ContentCard.vue';

import {usePullAirTemperature, useAirTemperature} from '@/traits/airTemperature/airTemperature.js';
import {computed} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: null
  },
  externalName: {
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

const {value: indoorValue} = usePullAirTemperature(() => props.name);
const {
  temp: indoorTemperature,
  humidity: indoorHumidity,
  tempRange
} = useAirTemperature(indoorValue);
const {value: outdoorValue} = usePullAirTemperature(() => props.externalName);
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
