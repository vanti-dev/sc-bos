<template>
  <div>
    <slot :resource="enterLeaveEventValue"/>
  </div>
</template>

<script setup>
import {usePullEnterLeaveEvents} from '@/traits/enterLeave/enterLeave.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullEnterLeaveEventsRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const enterLeaveEventValue = reactive(usePullEnterLeaveEvents(() => props.request ?? props.name, () => props.paused));
</script>
