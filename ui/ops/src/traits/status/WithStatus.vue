<template>
  <div>
    <slot :resource="statusLogValue"/>
  </div>
</template>

<script setup>
import {usePullCurrentStatus} from '@/traits/status/status.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullCurrentStatusRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const statusLogValue = reactive(usePullCurrentStatus(() => props.request ?? props.name, () => props.paused));
</script>
