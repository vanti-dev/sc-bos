<template>
  <div>
    <slot :resource="airQualityResource"/>
  </div>
</template>

<script setup>
import {usePullAirQuality} from '@/traits/airQuality/airQuality.js';
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

const airQualityResource = reactive(usePullAirQuality(() => props.name, () => props.paused));
</script>

<style lang="scss">
</style>
