<template>
  <div>
    <slot
        :occupant-count="occupantCount"
        :occupancy-state="occupancyState"
        :occupancy-value="occupancyValue"/>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {newResourceValue} from '@/api/resource';
import {occupancyStateToString} from '@/api/sc/traits/occupancy';

import {useTableDataStore} from '@/stores/tableDataStore';
import {useErrorStore} from '@/components/ui-error/error';

const {handleStream} = useTableDataStore();
const errorStore = useErrorStore();

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  paused: Boolean
});

const occupancyValue = reactive(
    /** @type {ResourceValue<Occupancy.AsObject, Occupancy>} */
    newResourceValue()
);

//
//
// Computed
const occupantCount = computed(() => {
  if (occupancyValue.value) {
    return occupancyValue.value.peopleCount;
  }
  return 0;
});

const occupancyState = computed(() => {
  if (occupancyValue.value) {
    return occupancyStateToString(occupancyValue.value.state);
  }
  return 'unknown';
});


//
//
// Watch
watch(() => ([props.name, props.paused]), () => {
  // pinia action
  handleStream(props.name, props.paused, occupancyValue);
}, {immediate: true});


//
//
// UI error handling
let unwatchOccupancyError;

onMounted(() => {
  unwatchOccupancyError = errorStore.registerValue(occupancyValue);
});

onUnmounted(() => {
  if (unwatchOccupancyError) unwatchOccupancyError();
});
</script>

<style lang="scss">
.occupied {
  color: var(--v-success-lighten1) !important;
}
.idle {
  color: var(--v-info-base) !important;
}
.unoccupied {
  color: var(--v-warning-base) !important;
}
</style>
