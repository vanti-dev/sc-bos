<template>
  <content-card class="mb-5 d-flex flex-column pt-7 pb-0">
    <h4 class="text-h4 pl-4 pb-8 pt-0">Environmental</h4>
    <div class="d-flex flex-row flex-nowrap">
      <circular-gauge
          :value="temperature"
          color="#ffc432"
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
          color="#ffc432"
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
  const oldNames = Array.isArray(oldName) ? oldName : [oldName];
  const newNames = Array.isArray(newName) ? newName : [newName];

  // Remove old names
  oldNames.forEach((name) => {
    if (!newNames.includes(name)) {
      unwatchErrorFunctions[name]?.(); // unwatch the error function
      delete unwatchErrorFunctions[name]; // delete the error function
      delete indoorValues[name]; // delete the value
    }
  });

  // Add new names
  newNames.forEach((name) => {
    if (!oldNames.includes(name) && name) {
      indoorValues[name] = useAirTemperatureTrait({name, paused: false});
      // watch the error function
      unwatchErrorFunctions[name] = errorStore.registerValue(indoorValues[name].airTemperatureResource);
    }
  });
}, {immediate: true, deep: true});


// Watch for changes to the externalName prop and update the outdoorValues object
watch(() => props.externalName, (newName, oldName) => {
  const oldNames = Array.isArray(oldName) ? oldName : [oldName];
  const newNames = Array.isArray(newName) ? newName : [newName];

  // Remove old names
  oldNames.forEach((name) => {
    if (!newNames.includes(name)) {
      unwatchErrorFunctions[name]?.(); // unwatch the error function
      delete unwatchErrorFunctions[name]; // delete the error function
      delete outdoorValues[name]; // delete the value
    }
  });

  // Add new names
  newNames.forEach((name) => {
    if (!oldNames.includes(name) && name) {
      outdoorValues[name] = useAirTemperatureTrait({name, paused: false});
      // watch the error function
      unwatchErrorFunctions[name] = errorStore.registerValue(outdoorValues[name].airTemperatureResource);
    }
  });
}, {immediate: true, deep: true});


// ------------------------------------ //
// Calculate the average indoor temperature (if multiple devices are specified)
// or return the temperature of the single device specified
const averageIndoorTempValue = computed(() => {
  if (Array.isArray(props.name)) {
    const values = props.name.map(
        (name) => indoorValues[name].airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0);
    const sum = values.reduce((a, b) => a + b, 0);
    return sum / values.length;
  } else {
    return indoorValues[props.name]?.airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0;
  }
});

// Calculate the average indoor humidity (if multiple devices are specified)
// or return the humidity of the single device specified
const averageIndoorHumidityValue = computed(() => {
  if (Array.isArray(props.name)) {
    const values = props.name.map(
        (name) => indoorValues[name].airTemperatureResource?.value?.ambientHumidity ?? 0);
    const sum = values.reduce((a, b) => a + b, 0);
    return sum / values.length;
  } else {
    return indoorValues[props.name]?.airTemperatureResource?.value?.ambientHumidity ?? 0;
  }
});

// Calculate the average external temperature (if multiple external devices are specified)
// or return the external temperature of the single device specified
const averageOutdoorTempValue = computed(() => {
  if (Array.isArray(props.externalName)) {
    const values = props.externalName.map(
        (name) => outdoorValues[name].airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0);
    const sum = values.reduce((a, b) => a + b, 0);
    return sum / values.length;
  } else {
    return outdoorValues[props.externalName]?.airTemperatureResource?.value?.ambientTemperature?.valueCelsius ?? 0;
  }
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
