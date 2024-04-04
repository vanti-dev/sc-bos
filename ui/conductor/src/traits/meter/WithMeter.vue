<template>
  <div>
    <slot :resource="meterReadings" :info="meterReadingInfo"/>
  </div>
</template>

<script setup>
import {useDescribeMeterReading, usePullMeterReading} from '@/traits/meter/meter.js';
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

const meterReadings = reactive(usePullMeterReading(() => props.name, () => props.paused));
const meterReadingInfo = reactive(useDescribeMeterReading(() => props.name));
</script>
