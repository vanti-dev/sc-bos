<template>
  <span class="d-flex flex-row flex-nowrap">
    <span>
      <template v-if="showErr">
        <v-tooltip location="bottom">
          <template #activator="{props: _props}">
            <v-icon v-bind="_props" color="error" size="1em">mdi-alert-circle-outline</v-icon>
          </template>
          <span>{{ errStr }}</span>
        </v-tooltip>
      </template>
      <span v-else class="value">{{ props.peopleCount }}</span>
      <template v-if="maxOccupancy > 0">
        <span class="div">/</span>
        <span class="total">{{ props.maxOccupancy }}</span>
      </template>
    </span>
    <span v-if="maxOccupancy > 0" class="ml-5 text-right" style="min-width: 2.5em">
      <span class="value">{{ occupancyPercentageDisplay }}</span>
      <span class="unit">%</span>
    </span>
  </span>
</template>

<script setup>
import useError from '@/composables/error.js';
import {computed} from 'vue';

const props = defineProps({
  peopleCount: {
    type: Number,
    default: 0
  },
  maxOccupancy: {
    type: Number,
    default: 1625
  },
  error: {
    type: [Object, String],
    default: null
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

const {errStr, showErr} = useError(() => props.error);
</script>

<style scoped>
.div, .total, .unit {
  font-size: 50%;
  opacity: 0.8;
  font-weight: lighter;
  margin-left: 0.2em;
}
</style>
