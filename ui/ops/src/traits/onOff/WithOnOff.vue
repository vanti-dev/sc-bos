<template>
  <div>
    <slot :resource="onOffValue" :update="doUpdateOnOff" :update-tracker="updateTracker"/>
  </div>
</template>

<script setup>
import {usePullOnOff, useUpdateOnOff} from '@/traits/onOff/onOff.js';
import {reactive} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const onOffValue = reactive(usePullOnOff(() => props.name, () => props.paused));
const updateTracker = reactive(useUpdateOnOff(() => props.name));
const doUpdateOnOff = updateTracker.updateOnOff;
</script>
