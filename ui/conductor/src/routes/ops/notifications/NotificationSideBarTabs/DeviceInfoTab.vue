<template>
  <span>
    <device-side-bar-content :device-id="deviceId" :traits="traits"/>
  </span>
</template>

<script setup>
import DeviceSideBarContent from '@/routes/devices/components/DeviceSideBarContent.vue';
import {useIntersectedItemsStore} from '@/stores/intersectedItemsStore';
import {computed, watch} from 'vue';

const {intersectionHandler} = useIntersectedItemsStore();

const props = defineProps({
  deviceId: {
    type: String,
    default: ''
  },
  deviceData: {
    type: Object,
    default: () => {
    }
  }
});

const traits = computed(() => {
  const traits = {};
  if (props.deviceData?.metadata?.traitsList) {
    // flatten array of trait objects (e.g. [{name: 'trait1', ...}, ...] into object (e.g. {trait1: true, ...})
    props.deviceData.metadata.traitsList.forEach((trait) => traits[trait.name] = true);
  }

  return traits;
});

// Update the intersectionHandler when the deviceId changes
// so we can establish stream and pull in the device data
watch(() => props.deviceId, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    if (!newVal) intersectionHandler([{isIntersecting: false}], null, oldVal);
    else {
      intersectionHandler([{isIntersecting: true}], null, newVal);
      intersectionHandler([{isIntersecting: false}], null, oldVal);
    }
  }
}, {immediate: true, deep: true});
</script>
