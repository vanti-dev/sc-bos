<template>
  <div>
    <slot :resource="occupancyValue"/>
  </div>
</template>

<script setup>
import {usePullOccupancy} from '@/traits/occupancy/occupancy.js';
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

const occupancyValue = reactive(usePullOccupancy(() => props.name, () => props.paused));
</script>
