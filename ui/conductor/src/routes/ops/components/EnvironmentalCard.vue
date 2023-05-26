<template>
  <content-card class="mb-8 d-flex flex-column px-6 pt-md-6">
    <h4 class="text-h4 order-lg-last pb-4 pb-lg-0 pt-0 pt-lg-4">Environmental</h4>
    <div class="d-flex flex-column flex-md-row flex-lg-column">
      <circular-gauge
          :value="temperature"
          :min="temperatureRange.low"
          :max="temperatureRange.high"
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
          segments="30"
          class="align-self-center">
        <span class="align-baseline">
          {{ (humidity * 100).toFixed(1) }}<span style="font-size: 0.7em;">&percnt;</span>
        </span>
        <template #title>Avg. Humidity</template>
      </circular-gauge>
    </div>
  </content-card>
</template>

<script setup>

import ContentCard from '@/components/ContentCard.vue';
import CircularGauge from '@/components/CircularGauge.vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullAirTemperature} from '@/api/sc/traits/air-temperature';
import {useErrorStore} from '@/components/ui-error/error';

const props = defineProps({
  // name of the device/zone to query for internal temperature data
  name: {
    type: String,
    default: 'building'
  },
  // name of the device/zone to query for external temperature data
  externalName: {
    type: String,
    default: 'outside'
  }
});

// todo: do we need to get this from somewhere?
const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});

const indoorTempValue = reactive(/** @type {ResourceValue<AirTemperature.AsObject, AirTemperature>} */ newResourceValue());
const outdoorTempValue = reactive(/** @type {ResourceValue<AirTemperature.AsObject, AirTemperature>} */ newResourceValue());

const temperature = computed(() => indoorTempValue.value?.ambientTemperature?.valueCelsius ?? 0);
const humidity = computed(() => indoorTempValue.value?.ambientHumidity ?? 0);
const externalTemperature = computed(() => outdoorTempValue.value?.ambientTemperature?.valueCelsius ?? 0);

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(indoorTempValue);
  // create new stream
  if (name && name !== '') {
    pullAirTemperature(name, indoorTempValue);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(indoorTempValue);
});

watch(() => props.externalName, async (name) => {
  // close existing stream if present
  closeResource(outdoorTempValue);
  // create new stream
  if (name && name !== '') {
    pullAirTemperature(name, outdoorTempValue);
  }
}, {immediate: true});

onUnmounted(() => {
  closeResource(outdoorTempValue);
});

// UI Error handling
const errorStore = useErrorStore();
let unwatchIndoorTempErrors; let unwatchOutdoorTempErrors;
onMounted(() => {
  unwatchIndoorTempErrors = errorStore.registerValue(indoorTempValue);
  unwatchOutdoorTempErrors = errorStore.registerValue(outdoorTempValue);
});
onUnmounted(() => {
  unwatchIndoorTempErrors();
  unwatchOutdoorTempErrors();
});


</script>

<style scoped>

</style>
