<template>
  <div>
    <slot :resource="emergencyValue"/>
  </div>
</template>

<script setup>
import {usePullEmergency} from '@/traits/emergency/emergency.js';
import {reactive} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullEmergencyRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const emergencyValue = reactive(usePullEmergency(() => props.request ?? props.name, () => props.paused));
</script>
