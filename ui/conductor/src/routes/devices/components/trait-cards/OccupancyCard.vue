<template>
  <v-card elevation="0" tile>
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Occupancy Sensor</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">State</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize" :class="state.toLowerCase()">{{ state }}</v-list-item-subtitle>
      </v-list-item>
      <v-list-item class="py-1" v-if="count !== 0">
        <v-list-item-title class="text-body-small text-capitalize">Count</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ count }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear color="primary" indeterminate :active="occupancyValue.loading"/>
  </v-card>
</template>

<script setup>

import {closeResource, newResourceValue} from '@/api/resource';
import {occupancyStateToString, pullOccupancy} from '@/api/sc/traits/occupancy';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  }
});

const occupancyValue = reactive(/** @type{ResourceValue<Occupancy.AsObject, Occupancy>} */newResourceValue());

const state = computed(() => {
  if (occupancyValue.value) {
    return occupancyStateToString(occupancyValue.value.state);
  }
  return 'unknown';
});
const count = computed(() => {
  if (occupancyValue.value) {
    return occupancyValue.value.peopleCount;
  }
  return 0;
});

watch(() => props.name, async (name) => {
  // close existing stream if present
  closeResource(occupancyValue);
  // create new stream
  if (name && name !== '') {
    await pullOccupancy(name, occupancyValue);
  }
}, {immediate: true});

// UI error handling
const errorStore = useErrorStore();
let unwatchOccupancyError;
onMounted(() => {
  unwatchOccupancyError = errorStore.registerValue(occupancyValue);
});
onUnmounted(() => {
  if (unwatchOccupancyError) unwatchOccupancyError();
});

</script>

<style scoped>
.v-list-item {
    min-height: auto;
}
.v-list-item__subtitle.occupied {
    color: var(--v-success-lighten1) !important;
}
.v-list-item__subtitle.idle {
    color: var(--v-info-base) !important;
}
.v-list-item__subtitle.unoccupied {
    color: var(--v-warning-base) !important;
}
</style>
