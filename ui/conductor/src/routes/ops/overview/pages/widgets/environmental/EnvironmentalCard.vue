<template>
  <content-card class="mb-5 d-flex flex-column pt-7 pb-0">
    <h4 class="text-h4 pl-4 pb-8 pt-0">Environmental</h4>
    <div :class="['d-flex flex-row justify-center', {'flex-wrap mb-4': props.shouldWrap}]">
      <circular-gauge
          :value="temperature"
          :color="props.gaugeColor"
          :min="temperatureRange.low"
          :max="temperatureRange.high"
          segments="30"
          style="max-width: 140px;"
          class="mt-2 mb-5 ml-3 mr-2">
        <span class="mt-n4 ml-1 text-h1">
          {{ temperature.toFixed(1) }}&deg;
        </span>
        <template #title>
          <span class="ml-n1 mb-2">Avg. Indoor Temperature</span>
        </template>
      </circular-gauge>
      <div
          v-if="externalName"
          :class="[humidity > 0 ? 'mb-7' : 'mb-2',
                   'd-flex flex-column justify-end align-center']"
          style="width: 150px;">
        <span
            class="text-h1 align-left mb-3"
            style="display: inline-block;">{{ externalTemperature.toFixed(1) }}&deg;
        </span>
        <span
            class="text-title text-center"
            style="display: inline-block; width: 100px;">
          External Temperature
        </span>
      </div>
      <circular-gauge
          v-if="humidity > 0"
          :value="humidity"
          :color="props.gaugeColor"
          segments="30"
          style="max-width: 140px;"
          class="mt-2">
        <span class="align-baseline text-h1 mt-n2">
          {{ (humidity * 100).toFixed(1) }}<span style="font-size: 0.7em;">%</span>
        </span>
        <template #title>
          <span class="mb-2">Avg. Humidity</span>
        </template>
      </circular-gauge>
    </div>
  </content-card>
</template>

<script setup>
import CircularGauge from '@/components/CircularGauge.vue';
import ContentCard from '@/components/ContentCard.vue';

import {useErrorStore} from '@/components/ui-error/error';
import useAirTemperatureTrait from '@/composables/traits/useAirTemperatureTrait';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';

const props = defineProps({
  // name of the device/zone to query for internal temperature data
  name: {
    type: String,
    default: ''
  },
  // name of the device/zone to query for external temperature data
  externalName: {
    type: String,
    default: ''
  },
  gaugeColor: {
    type: String,
    default: ''
  },
  shouldWrap: {
    type: Boolean,
    default: false
  }
});

// todo: do we need to get this from somewhere?
const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});
const indoorValues = reactive({});
const outdoorValues = reactive({});


// Error handling
const errorStore = useErrorStore();
const unwatchErrorFunctions = {};


// Watch for changes to the name prop and update the indoorValues object
watch(() => props.name, (newName, oldName) => {
  // Remove old name
  useAirTemperatureTrait({name: oldName, paused: true}).clearResourceError();
  unwatchErrorFunctions[oldName]?.(); // unwatch the error function
  delete unwatchErrorFunctions[oldName]; // delete the error function
  delete indoorValues[oldName]; // delete the value

  if (newName) {
    // Add new names
    indoorValues[newName] = useAirTemperatureTrait({name: newName, paused: false});
    // watch the error function
    unwatchErrorFunctions[newName] = errorStore.registerValue(indoorValues[newName].airTemperatureResource);
  }
}, {immediate: true, deep: true});


// Watch for changes to the externalName prop and update the outdoorValues object
watch(() => props.externalName, (newName, oldName) => {
  // Remove old name
  useAirTemperatureTrait({name: oldName, paused: true}).clearResourceError();
  unwatchErrorFunctions[oldName]?.(); // unwatch the error function
  delete unwatchErrorFunctions[oldName]; // delete the error function
  delete outdoorValues[oldName]; // delete the value


  if (newName) {
    // Add new name
    outdoorValues[newName] = useAirTemperatureTrait({name: newName, paused: false});
    // watch the error function
    unwatchErrorFunctions[newName] = errorStore.registerValue(outdoorValues[newName].airTemperatureResource);
  }
}, {immediate: true, deep: true});


// ------------------------------------ //
// Return the temperature of the single device specified
const averageIndoorTempValue = computed(() => {
  return indoorValues[props.name]?.airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0;
});

// Return the humidity of the single device specified
const averageIndoorHumidityValue = computed(() => {
  return indoorValues[props.name]?.airTemperatureResource?.value?.ambientHumidity ?? 0;
});

// Return the external temperature of the single device specified
const averageOutdoorTempValue = computed(() => {
  return outdoorValues[props.externalName]?.airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0;
});

const temperature = computed(() => {
  return averageIndoorTempValue.value;
});
const humidity = computed(() => {
  return averageIndoorHumidityValue.value;
});

const externalTemperature = computed(() => {
  return averageOutdoorTempValue.value;
}
);

// ------------------------------------ //
// Clean up UI Error handling
onUnmounted(() => {
  Object.values(unwatchErrorFunctions).forEach(unwatch => unwatch());
});
</script>
