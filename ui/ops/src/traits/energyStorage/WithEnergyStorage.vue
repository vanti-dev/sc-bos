<template>
  <div>
    <slot :resource="energyValue"/>
  </div>
</template>

<script setup>
import {usePullEnergyLevel} from '@/traits/energyStorage/energyStorage.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullEnergyLevelRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const energyValue = reactive(usePullEnergyLevel(() => props.request ?? props.name, () => props.paused));
</script>

