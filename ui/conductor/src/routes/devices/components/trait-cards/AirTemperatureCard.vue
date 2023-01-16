<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Temperature</v-subheader>
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
    <v-card-actions class="px-4">
      <v-spacer/>
      <v-btn small color="neutral lighten-1" elevation="0" @click="changeSetPoint(0.1)">Up</v-btn>
      <v-btn small color="neutral lighten-1" elevation="0" @click="changeSetPoint(-0.1)">Down</v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="updateValue.loading"/>
  </v-card>
</template>

<script setup>
import {computed, onUnmounted, reactive, ref, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {pullAirTemperature, toDisplayObject, updateAirTemperature} from '@/api/sc/traits/air-temperature';
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
}, {immediate: true});

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

const updateValue = reactive(newResourceValue());

/**
 * @param {number} value
 */
function changeSetPoint(value) {
  if (airTempValue.value &&
      airTempValue.value.hasOwnProperty('temperatureSetPoint') &&
      airTempValue.value.temperatureSetPoint !== undefined &&
      airTempValue.value.temperatureSetPoint.hasOwnProperty('valueCelsius')) {
    /* @type {UpdateAirTemperatureRequest.AsObject} */
    const req = {
      name: props.name,
      state: {
        temperatureSetPoint: {
          valueCelsius: airTempValue.value.temperatureSetPoint.valueCelsius + value
        }
      }
      // todo: add updateMask?
    };
    updateAirTemperature(req, updateValue);
  }
}

</script>

<style scoped>
.v-list-item {
  min-height: auto;
}
.v-progress-linear {
  width: auto;
}
</style>
