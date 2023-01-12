<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small  text-capitalize">Set Point</v-list-item-title>
        <v-list-item-subtitle>{{ latestTemperature.temperatureSetPoint.val }}°C</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Current Temp</v-list-item-title>
        <v-list-item-subtitle>{{ latestTemperature.ambientTemperature.val }}°C</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="tempValue()"
        background-color="neutral lighten-1"
        color="accent"/>
  </v-card>
</template>

<script setup>
import {defineProps, ref} from 'vue';

const temperatureRange = ref({
  low: 18.0,
  high: 24.0
});

const latestTemperature = ref({
  ambientTemperature: {
    val: 21.5
  },
  temperatureSetPoint: {
    val: 22.0
  }
});

defineProps({
  airTemperaturePull: {
    type: Object,
    default: () => ({
      mode: 'MODE_UNSPECIFIED',
      ambient_temperature: null,
      dew_point: null
    })
  }
});

/**
 *  @return {number}
 */
function tempValue() {
  let val = latestTemperature.value.ambientTemperature.val;
  val -= temperatureRange.value.low;
  val = val/(temperatureRange.value.high-temperatureRange.value.low);
  console.log(val);
  return val*100;
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
