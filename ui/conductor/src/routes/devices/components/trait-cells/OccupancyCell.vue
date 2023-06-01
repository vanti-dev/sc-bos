<template>
  <v-tooltip left transition="slide-x-reverse-transition" :color="colorStr">
    <template #activator="{on}">
      <v-icon :class="state" :color="iconColor" v-on="on" size="20">{{ iconStr }}</v-icon>
    </template>
    <span>{{ stateStr }}</span>
  </v-tooltip>
</template>

<script setup>
import {occupancyStateToString} from '@/api/sc/traits/occupancy';
import {Occupancy} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {
    }
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const state = computed(() => {
  return props.value?.state;
});
const stateStr = computed(() => {
  if (state.value === undefined) return '';
  return occupancyStateToString(state.value);
});
const colorStr = computed(() => {
  return stateStr.value.toLowerCase();
});

const iconStr = computed(() => {
  if (state.value === Occupancy.State.OCCUPIED) {
    return 'mdi-crosshairs-gps';
  } else if (state.value === Occupancy.State.UNOCCUPIED) {
    return 'mdi-crosshairs';
  } else if (state.value === Occupancy.State.IDLE) {
    return 'mdi-crosshairs-gps';
  } else {
    return '';
  }
});
const iconColor = computed(() => {
  if (state.value === Occupancy.State.IDLE) {
    return 'grey';
  } else {
    return undefined;
  }
});
</script>
