<template>
  <div>
    <slot :resource="onOffValue"/>
  </div>
</template>

<script setup>
import {usePullOnOff} from '@/traits/onOff/onOff.js';
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
</script>
