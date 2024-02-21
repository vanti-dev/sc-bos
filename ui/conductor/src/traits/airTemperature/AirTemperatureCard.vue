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
          @click="changeSetPoint(-0.5)"
          :disabled="blockActions || (props.value?.temperatureSetPoint === undefined)">
        Down
      </v-btn>
      <v-btn
          small
          color="neutral lighten-1"
          elevation="0"
          @click="changeSetPoint(0.5)"
          :disabled="blockActions || (props.value?.temperatureSetPoint === undefined)">
        Up
      </v-btn>
    </v-card-actions>
    <v-progress-linear color="primary" indeterminate :active="updateValue.loading"/>
  </v-card>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import {airTemperatureModeToString, temperatureToString} from '@/api/sc/traits/air-temperature';
import useAuthSetup from '@/composables/useAuthSetup';
import {camelToSentence} from '@/util/string';
import {AirTemperature} from '@smart-core-os/sc-api-grpc-web/traits/air_temperature_pb';
import {computed, reactive, ref} from 'vue';

const {blockActions} = useAuthSetup();


const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});

const props = defineProps({
  value: {
    type: Object, // of AirTemperature.AsObject
    default: () => ({})
  },
  loading: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits([
  'updateAirTemperature' // of number | AirTemperature.AsObject | UpdateAirTemperatureRequest.AsObject
]);

const updateValue = reactive(newActionTracker());

/**
 * Calculates the percentage value of the current temperature based on the temperature range
 *
 * @return {number}
 */
function tempProgress() {
  let val = props.value?.ambientTemperature?.valueCelsius ?? 0;
  if (val > 0) {
    val -= temperatureRange.value.low;
    val = val / (temperatureRange.value.high - temperatureRange.value.low);
  }
  return val * 100;
}

const airTempData = computed(() => {
  if (props && props.value) {
    const data = {};
    Object.entries(props.value).forEach(([key, value]) => {
      if (value !== undefined) {
        switch (key) {
          case 'mode':
            if (value !== AirTemperature.Mode.MODE_UNSPECIFIED) {
              data[key] = airTemperatureModeToString(value);
            }
            break;
          case 'ambientTemperature': {
            data['currentTemp'] = temperatureToString(value);
            break;
          }
          case 'temperatureSetPoint': {
            data['setPoint'] = temperatureToString(value);
            break;
          }
          case 'ambientHumidity':
            if (value !== 0) {
              data['humidity'] = value.toFixed(1) + '%';
            }
            break;
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
  if (props.value?.temperatureSetPoint?.valueCelsius !== undefined) {
    emit('updateAirTemperature', props.value.temperatureSetPoint.valueCelsius + value);
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
