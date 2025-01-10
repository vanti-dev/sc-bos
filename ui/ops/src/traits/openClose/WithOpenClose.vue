<template>
  <div>
    <slot :resource="openCloseValue"/>
  </div>
</template>

<script setup>
import {usePullOpenClosePositions} from '@/traits/openClose/openClose.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullAccessAttemptsRequest.AsObject
    default: () => {}
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const openCloseValue = reactive(usePullOpenClosePositions(() => props.request ?? props.name, () => props.paused));
</script>
