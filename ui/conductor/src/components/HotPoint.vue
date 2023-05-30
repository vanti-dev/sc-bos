<template>
  <div v-intersect="onIntersect">
    <slot name="hotpoint" :live="isLive" :sensor-types="sensorTypes"/>
  </div>
</template>

<script setup>
import {computed, onBeforeUnmount, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';

import {useIntersectedItemsStore} from '@/stores/intersectedItemsStore';

import {Device} from '@sc-bos/ui-gen/proto/devices_pb';

const props = defineProps({
  item: {
    type: Device,
    default: () => {}
  },
  itemKey: {
    type: String,
    default: ''
  }
});
const intersectedItemsStore = useIntersectedItemsStore();
const {clearName, createName, intersectionHandler} = intersectedItemsStore;
const {intersectedItemNames} = storeToRefs(intersectedItemsStore);
const isLive = ref(true);

const onIntersect = {
  handler: (
      entries,
      observer
  ) =>
    intersectionHandler(entries, observer, props.itemKey),
  options: {
    rootMargin: '-50px 0px 0px 0px',
    threshold: 0.75,
    trackVisibility: true,
    delay: 100
  }
};

//
//
// Computed
// Return all trait (sensor) type
const sensorTypes = computed(() => {
  const traitsArray = props.item?.metadata?.traitsList.map(trait => {
    return trait.name;
  });

  const sensors = traitsArray.map((item) => {
    const arr = item.split('.');
    return arr.slice(2).join('.');
  });

  return sensors;
});

//
//
// Watchers
// Watching (matching) the live/paused state
watch(() => intersectedItemNames.value, names => {
  if (names[props.itemKey]) isLive.value = true;
  else isLive.value = false;
}, {immediate: true, deep: true, flush: 'sync'});

// Updating intersectedItemNames if there is a device name change
watch(() => props.itemKey, (newKey, oldKey) => {
  if (newKey !== oldKey) {
    clearName(oldKey);
    createName(newKey);
  }
}, {immediate: true, deep: true, flush: 'sync'});

//
//
// Lifecycle
onBeforeUnmount(() => {
  clearName(props.itemKey);
});
</script>
