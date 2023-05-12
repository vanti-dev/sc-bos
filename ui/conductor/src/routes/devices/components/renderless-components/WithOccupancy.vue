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
import {closeResource, newResourceValue} from '@/api/resource';
import {occupancyStateToString, pullOccupancy} from '@/api/sc/traits/occupancy';

import {useErrorStore} from '@/components/ui-error/error';
import {usePageStore} from '@/stores/page';
import {Device} from '@sc-bos/ui-gen/proto/devices_pb';
const errorStore = useErrorStore();
const pageStore = usePageStore();

const props = defineProps({
  item: {
    type: Device,
    default: () => {}
  },
  itemName: {
    type: String,
    default: ''
  },
  table: {
    type: Boolean,
    default: true
  }
});

const occupancyValue = reactive(/** @type{ResourceValue<Occupancy.AsObject, Occupancy>} */newResourceValue());

//
//
// Computed
const occupancyState = computed(() => {
  if (occupancyValue.value) {
    return occupancyStateToString(occupancyValue.value.state);
  }
  return 'unknown';
});

const occupantCount = computed(() => {
  if (occupancyValue.value) {
    return occupancyValue.value.peopleCount;
  }
  return 0;
});

//
//
// Watchers

// If we looking at the table only
if (props.table) {
  // Let's see if the row is intersecting
  watch(() => props.item.isIntersected, (newValue, oldValue) => {
    // console.log('Value changed from', oldValue, 'to', newValue);

    // If the row is NOT intersecting, close any existing connection
    if (!oldValue || !newValue) closeResource(occupancyValue);

    // name = name value from Table
    const name = props.item.name;

    // If the row is intersecting, open a new connection
    if (newValue) pullOccupancy(name, occupancyValue);
    // called immediately for only the desired value, asynchronously after all reactive updates
  }, {immediate: true, deep: false, flush: 'post'});
  //
  // Other than the table
} else {
  // Let's watch the item's name
  watch(() => ({
    name: props.itemName,
    showSidebar: pageStore.showSidebar
  }), (newValues, oldValues) => {
    // console.log('Name changed from', oldValues?.name, 'to', newValues?.name);

    // If the sidebar visible
    if (newValues.showSidebar) {
      // then if there is a name change
      if (
        newValues?.name !== oldValues?.name ||
        newValues.name === '' ||
        !newValues.name
      ) closeResource(occupancyValue);

      // If we actually have a name (selected item), then make a call
      if (newValues.name !== '') pullOccupancy(newValues.name, occupancyValue);
      //
      // finally if the sidebar hidden
    } else if (!newValues.showSidebar) {
      // If we have a new name or
      closeResource(occupancyValue);
      // console.log('Connection closed');
    }
    // called immediately for only the desired value, synchronously immediately after the watched value changes
  }, {immediate: true, deep: true, flush: 'sync'});
}

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
