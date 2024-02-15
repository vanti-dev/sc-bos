<template>
  <content-card class="mb-5 d-flex flex-column pt-7 pb-0">
    <h4 class="text-h4 pl-4 pb-8 pt-1">Environmental</h4>
    <div class="d-flex flex-column align-center mb-4">
      <v-col cols="auto" class="ma-0 pa-0">
        <circular-gauge
            v-if="indoorTemperature > 0 || props.shouldWrap"
            :value="indoorTemperature"
            :color="props.gaugeColor"
            :min="temperatureRange.low"
            :max="temperatureRange.high"
            segments="30"
            style="max-width: 140px;"
            class="mt-2 mb-5 ml-3 mr-2">
          <span class="mt-n4 ml-1 text-h1">
            {{ indoorTemperature.toFixed(1) }}&deg;
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
                     'd-flex flex-column justify-end align-center']"
            style="width: 150px;">
          <span
              class="text-h1 align-left mb-3"
              style="display: inline-block;">{{ outdoorTemperature.toFixed(1) }}&deg;
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
            segments="30"
            style="max-width: 140px;"
            class="mt-2">
          <span class="align-baseline text-h1 mt-n2">
            {{ (indoorHumidity * 100).toFixed(1) }}<span style="font-size: 0.7em;">%</span>
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

import useAirTemperatureTrait from '@/composables/traits/useAirTemperatureTrait';
import {computed, reactive, ref, watchEffect} from 'vue';

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

const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});

const indoorProps = reactive({
  name: props.name,
  paused: false
});
const outdoorProps = reactive({
  name: props.externalName,
  paused: false
});

watchEffect(() => {
  indoorProps.name = props.name;
  outdoorProps.name = props.externalName;
});

const indoor = useAirTemperatureTrait(indoorProps);
const outdoor = useAirTemperatureTrait(outdoorProps);

const indoorTemperature = computed(() => {
  return indoor.temperatureValue.value;
});

const outdoorTemperature = computed(() => {
  return outdoor.temperatureValue.value;
});

const indoorHumidity = computed(() => {
  return indoor.humidityValue.value;
});
</script>
