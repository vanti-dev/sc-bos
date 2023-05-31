<template>
  <div v-intersect="onIntersect">
    <slot name="hotpoint" :live="isLive"/>
  </div>
</template>

<script setup>
import {useIntersectedItemsStore} from '@/stores/intersectedItemsStore';
import {storeToRefs} from 'pinia';
import {onBeforeUnmount, ref, watch} from 'vue';

const props = defineProps({
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
    // 60 for the page header
    rootMargin: '-60px 0px 0px 0px',
    threshold: 0
  }
};

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
  if (newKey !== oldKey && oldKey) {
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
