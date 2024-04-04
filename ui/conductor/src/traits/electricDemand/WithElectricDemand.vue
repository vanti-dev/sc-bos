<template>
  <div>
    <slot :resource="demandValue"/>
  </div>
</template>

<script setup>
import {usePullElectricDemand} from '@/traits/electricDemand/electric.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullDemandRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const demandValue = reactive(usePullElectricDemand(() => props.request ?? props.name, () => props.paused));
</script>
