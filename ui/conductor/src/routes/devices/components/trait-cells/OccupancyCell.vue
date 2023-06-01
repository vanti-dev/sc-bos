<template>
  <v-menu
      left
      bottom
      offset-y
      nudge-bottom="4px"
      nudge-right="4px"
      transition="slide-x-reverse-transition"
      open-on-hover>
    <template #activator="{on}">
      <v-icon :class="state" :color="iconColor" v-on="on" size="20">{{ iconStr }}</v-icon>
    </template>
    <v-card :color="colorStr">
      <v-card-text class="py-2">{{ tooltipStr }}</v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb';
import {occupancyStateToString} from '@/api/sc/traits/occupancy';
import {Occupancy} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {computed, onMounted, onUnmounted, ref} from 'vue';

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
  if (state.value === Occupancy.State.OCCUPIED) {
    return 'success lighten1';
  } else if (state.value === Occupancy.State.UNOCCUPIED) {
    return 'warning';
  } else if (state.value === Occupancy.State.IDLE) {
    return 'info';
  } else {
    return undefined;
  }
});

const nowHandle = ref(0);
const now = ref(new Date());
onUnmounted(() => clearInterval(nowHandle.value));
onMounted(() => {
  nowHandle.value = setInterval(() => {
    now.value = new Date();
  }, 1000);
});

const lastUpdateDate = computed(() => {
  return timestampToDate(props.value?.stateChangeTime);
});
const millisSinceLastUpdate = computed(() => {
  return now.value.getTime() - lastUpdateDate.value.getTime();
});
const showTimeSinceUpdate = computed(() => {
  return lastUpdateDate.value && millisSinceLastUpdate.value > 1000;
});
const sinceLastUpdateStr = computed(() => {
  if (!showTimeSinceUpdate.value) return '';
  const t = millisSinceLastUpdate.value;
  if (t > 1000 * 60 * 60 * 24) {
    const h = Math.floor(t / (1000 * 60 * 60 * 24));
    return `${h}d`;
  } else if (t > 1000 * 60 * 60) {
    const h = Math.floor(t / (1000 * 60 * 60));
    return `${h}h`;
  } else if (t > 1000 * 60) {
    const m = Math.floor(t / (1000 * 60));
    return `${m}m`;
  } else if (t > 1000) {
    const s = Math.floor(t / 1000);
    return `${s}s`;
  } else {
    return '';
  }
});

const tooltipStr = computed(() => {
  let s = stateStr.value;
  const sinceStr = sinceLastUpdateStr.value;
  if (sinceStr) {
    s += ` for ${sinceStr}`;
  }
  return s;
});

</script>
