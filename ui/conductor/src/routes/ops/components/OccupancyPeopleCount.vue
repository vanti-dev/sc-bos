<template>
  <div class="d-flex flex-row flex-nowrap">
    <div class="text-h2 d-flex flex-row flex-nowrap align-end">
      {{ props.peopleCount }}
      <span class="text-body-1 neutral--text text--lighten-5 pb-1 ml-1">/ {{ maxOccupancy }}</span>
    </div>
    <div class="text-h2 d-flex flex-row flex-nowrap align-end ml-6">
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
  }
});

// TODO: get this from somewhere
const maxOccupancy = 1234;

//
//
// Methods
const occupancyPercentage = computed(() => {
  return (props.peopleCount / maxOccupancy) * 100;
});

const occupancyPercentageDisplay = computed(() =>
  occupancyPercentage.value > 0 ? occupancyPercentage.value.toFixed(1) : occupancyPercentage.value.toFixed(0)
);
</script>
