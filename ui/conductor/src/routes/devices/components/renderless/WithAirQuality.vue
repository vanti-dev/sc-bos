<template>
  <div>
    <slot :resource="airQualityResource"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullAirQualitySensor} from '@/api/sc/traits/air-quality-sensor';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, reactive, watch} from 'vue';

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

const errorStore = useErrorStore();

const airQualityResource = reactive(
    /** @type {ResourceValue<AirQuality.AsObject, PullAirQualityResponse>} */
    newResourceValue());

//
//
// Watch
// Depending on paused state/device name, we close/open data stream(s)
watch(
    [() => props.paused, () => props.name],
    ([newPaused, newName], [oldPaused, oldName]) => {
      if (newPaused === oldPaused && newName === oldName) return;

      if (newPaused) {
        closeResource(airQualityResource);
      }

      if (!newPaused && (oldPaused || newName !== oldName)) {
        closeResource(airQualityResource);
        pullAirQualitySensor({name: newName}, airQualityResource);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

//
//
// UI error handling
const errorHandlers = [];
onMounted(() => {
});
onUnmounted(() => {
  closeResource(airQualityResource);
  errorHandlers.forEach(unwatch => unwatch());
});
</script>

<style lang="scss">
</style>
