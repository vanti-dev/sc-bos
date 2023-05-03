<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Temperature</v-subheader>
      <v-list-item v-for="(val, key) of airTempData" :key="key" class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">{{ camelToSentence(key) }}</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize font-weight-medium">{{ val }}</v-list-item-subtitle>
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
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="changeSetPoint(-0.1)"
          :disabled="(airTempValue.value?.temperatureSetPoint === undefined)">
        Down
      </v-btn>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="changeSetPoint(0.1)"
          :disabled="(airTempValue.value?.temperatureSetPoint === undefined)">
        Up
      </v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="updateValue.loading"/>
  </v-card>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {
  airTemperatureModeToString,
  pullAirTemperature,
  temperatureToString,
  updateAirTemperature
} from '@/api/sc/traits/air-temperature';
import {camelToSentence} from '@/util/string';
import {useErrorStore} from '@/components/ui-error/error';

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
const updateValue = reactive(newActionTracker());

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

// UI error handling
const errorStore = useErrorStore();
let unwatchAirTempErrors; let unwatchUpdateErrors;
onMounted(() => {
  unwatchAirTempErrors = errorStore.registerValue(airTempValue);
  unwatchUpdateErrors = errorStore.registerTracker(updateValue);
});
onUnmounted(() => {
  if (unwatchAirTempErrors) unwatchAirTempErrors();
  if (unwatchUpdateErrors) unwatchUpdateErrors();
});

/**
 * Calculates the percentage value of the current temperature based on the temperature range
 *
 * @return {number}
 */
function tempProgress() {
  let val = airTempValue.value?.ambientTemperature?.valueCelsius ?? 0;
  if (val > 0) {
    val -= temperatureRange.value.low;
    val = val / (temperatureRange.value.high - temperatureRange.value.low);
  }
  return val*100;
}

const airTempData = computed(() => {
  if (airTempValue && airTempValue.value) {
    const data = {};
    Object.entries(airTempValue.value).forEach(([key, value]) => {
      if (value !== undefined) {
        switch (key) {
          case 'mode': {
            data[key] = airTemperatureModeToString(value);
            break;
          }
          case 'ambientTemperature': {
            data['currentTemp'] = temperatureToString(value);
            break;
          }
          case 'temperatureSetPoint': {
            data['setPoint'] = temperatureToString(value);
            break;
          }
          case 'ambientHumidity': {
            data['humidity'] = (value * 100).toFixed(1) + '%';
            break;
          }
          case 'dewPoint': {
            data[key] = temperatureToString(value);
            break;
          }
          default: {
            data[key] = value;
          }
        }
      }
    });
    return data;
  }
  return {};
});


/**
 * @param {number} value
 */
function changeSetPoint(value) {
  if (airTempValue.value?.temperatureSetPoint?.valueCelsius !== undefined) {
    /* @type {UpdateAirTemperatureRequest.AsObject} */
    const req = {
      name: props.name,
      state: {
        temperatureSetPoint: {
          valueCelsius: airTempValue.value.temperatureSetPoint.valueCelsius + value
        }
      },
      updateMask: {pathsList: ['temperature_set_point']}
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
