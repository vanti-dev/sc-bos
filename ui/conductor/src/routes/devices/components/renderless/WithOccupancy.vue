<template>
  <div>
    <slot :resource="occupancyValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';
import {onUnmounted, reactive, watch} from 'vue';


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
        pullOccupancy({name: newName}, occupancyValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

onUnmounted(() => {
  closeResource(occupancyValue);
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
