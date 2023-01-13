<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps neutral--text text--lighten-3">State</v-subheader>
      <v-list-item v-for="(val, key) of airTempData" :key="key" class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">{{ camelToSentence(key) }}</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ val }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="tempProgress()"
        background-color="neutral lighten-1"
        color="accent"/>
  </v-card>
</template>

<script setup>
import {computed, defineProps, onUnmounted, reactive, ref, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullAirTemperature, toDisplayObject} from '@/api/sc/traits/air-temperature';
import {camelToSentence} from '@/util/string';

const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});

const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

const airTempValue = reactive(newResourceValue());

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(airTempValue);
  // create new stream
  if (name && name !== '') {
    pullAirTemperature(name, airTempValue);
  }
});

onUnmounted(() => {
  closeResource(airTempValue);
});

/**
 * Calculates the percentage value of the current temperature based on the temperature range
 *
 * @return {number}
 */
function tempProgress() {
  let val = 0;
  if (airTempValue.value &&
      airTempValue.value.hasOwnProperty('ambientTemperature') &&
      airTempValue.value.ambientTemperature !== undefined &&
      airTempValue.value.ambientTemperature.hasOwnProperty('valueCelsius')) {
    val = airTempValue.value.ambientTemperature.valueCelsius;
    val -= temperatureRange.value.low;
    val = val / (temperatureRange.value.high - temperatureRange.value.low);
  }
  return val*100;
}

const airTempData = computed(() => {
  if (airTempValue && airTempValue.value) {
    return toDisplayObject(airTempValue.value);
  }
  return {};
});

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
.v-progress-linear {
  width: auto;
}
</style>
