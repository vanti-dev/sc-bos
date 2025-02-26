<template>
  <div>
    <slot :resource="transportValue" :info="transportInfo"/>
  </div>
</template>

<script setup>
import {useDescribeTransport, usePullTransport} from '@/traits/transport/transport.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullTransportRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const transportValue = reactive(usePullTransport(() => props.name || props.request, () => props.paused));
const transportInfo = reactive(useDescribeTransport(() => props.name));
</script>
