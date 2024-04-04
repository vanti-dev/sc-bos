<template>
  <div>
    <slot :resource="accessAttemptValue"/>
  </div>
</template>

<script setup>
import {usePullAccessAttempts} from '@/traits/access/access.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullAccessAttemptsRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const accessAttemptValue = reactive(usePullAccessAttempts(() => props.name || props.request, () => props.paused));

</script>
