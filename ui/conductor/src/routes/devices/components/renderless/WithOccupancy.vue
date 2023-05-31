<template>
  <div>
    <slot
        name="occupancy"
        :occupancy-data="{
          occupantCount,
          occupancyState,
          occupancyValue
        }"/>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newResourceValue} from '@/api/resource';
import {occupancyStateToString, pullOccupancy} from '@/api/sc/traits/occupancy';
import {useErrorStore} from '@/components/ui-error/error';

const errorStore = useErrorStore();
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
  return '';
});

//
//
// Watch
// Depending on paused state/device name, we close/open data stream(s)
watch(
    [() => props.paused, () => props.name],
    ([newPaused, newName], [oldPaused, oldName]) => {
      // only for OccupancySensor
      if (newPaused === oldPaused && newName === oldName) return;

      if (newPaused) {
        closeResource(occupancyValue);
      }

      if (!newPaused && (oldPaused || newName !== oldName)) {
        closeResource(occupancyValue);
        pullOccupancy(newName, occupancyValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

//
//
// UI error handling
let unwatchOccupancyError;

onMounted(() => {
  unwatchOccupancyError = errorStore.registerValue(occupancyValue);
});

onUnmounted(() => {
  closeResource(occupancyValue);
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
