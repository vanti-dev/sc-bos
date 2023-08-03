<template>
  <div class="d-flex flex-row flex-nowrap">
    <div class="text-h2 d-flex flex-row flex-nowrap align-end">
      {{ props.peopleCount }}
      <span class="text-body-1 neutral--text text--lighten-5 pb-1 ml-1">/ {{ props.maxOccupancy }}</span>
    </div>
    <div class="text-h2 d-flex flex-row flex-nowrap align-end ml-8">
      {{ occupancyPercentageDisplay }}
      <span class="text-body-1 neutral--text text--lighten-9 pb-1 ml-1">%</span>
    </div>
  </div>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  peopleCount: {
    type: Number,
    default: 0
  },
  maxOccupancy: {
    type: Number,
    default: 1625
  }
});

//
//
// Methods
const occupancyPercentage = computed(() => {
  return (props.peopleCount / props.maxOccupancy) * 100;
});

const occupancyPercentageDisplay = computed(() =>
  occupancyPercentage.value > 0 ? occupancyPercentage.value.toFixed(1) : occupancyPercentage.value.toFixed(0)
);
</script>
