<template>
  <div>
    <slot
        :name="props.deviceType"
        :[`${slotData(props.deviceType).name}`]="slotData(props.deviceType).data"/>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {newResourceValue} from '@/api/resource';
import {occupancyStateToString} from '@/api/sc/traits/occupancy';

import {useTableDataStore} from '@/stores/tableDataStore';
import {useErrorStore} from '@/components/ui-error/error';
import {Device} from '@sc-bos/ui-gen/proto/devices_pb';

const {handleStream} = useTableDataStore();
const errorStore = useErrorStore();

const props = defineProps({
  deviceType: {
    type: String,
    default: ''
  },
  item: {
    type: Device,
    default: () => {}
  },
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
// Methods

// Defining slot data depending on sensor (device) type
const slotData = (sensorType) => {
  let data = {};
  let name = '';

  if (sensorType === 'occupancy') {
    name = sensorType + 'Data';
    data = {
      occupantCount: occupantCount.value,
      occupancyState: occupancyState.value,
      occupancyValue: occupancyValue
    };
  }

  return {name, data};
};

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
