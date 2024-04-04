<template>
  <div>
    <slot :resource="airTemperatureResource" :update="doUpdateAirTemperature" :update-tracker="updateTracker"/>
  </div>
</template>

<script setup>
import {usePullAirTemperature, useUpdateAirTemperature} from '@/traits/airTemperature/airTemperature.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const airTemperatureResource = reactive(usePullAirTemperature(() => props.name, () => props.paused));
const updateTracker = reactive(useUpdateAirTemperature(() => props.name));
const doUpdateAirTemperature = updateTracker.updateAirTemperature;
</script>

<style lang="scss">
.occupied {
  color: var(--v-success-lighten1) !important;
}

.idle {
  color: var(--v-info-base) !important;
}

.unoccupied {
  color: var(--v-warning-base) !important;
}
</style>
