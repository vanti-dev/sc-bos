<template>
  <floor-list :floors="props.floors">
    <template #floor="{floor}">
      <device-cell :item="deviceItemByLevel[floor.level]"/>
    </template>
  </floor-list>
</template>

<script setup>
/**
 * @typedef {Object} Floor
 * @property {number} level - 0 for ground, negative for basements (counting down), positive for upper floors (counting up). Defaults to len - i - 1.
 * @property {string} zoneName - smart core name for the floor zone, we will interrogate this to capture metadata and trait info
 */
import FloorList from '@/components/FloorList.vue';
import {useDevicesCollection} from '@/composables/devices.js';
import DeviceCell from '@/routes/devices/components/DeviceCell.vue';
import {computed} from 'vue';

const props = defineProps({
  floors: {
    type: Array, // of Floor
    required: true
  }
});

const devicesReq = computed(() => {
  return {
    query: {
      conditionsList: [
        {field: 'name', stringIn: {stringsList: props.floors.map(f => f.zoneName)}},
      ]
    }
  }
});
const deviceCollection = useDevicesCollection(devicesReq)
const floorByZoneName = computed(() => {
  const res = {};
  for (const floor of props.floors) {
    res[floor.zoneName] = floor;
  }
  return res;
});
const deviceItemByLevel = computed(() => {
  const byZoneName = floorByZoneName.value;
  const mdItems = deviceCollection.items.value;
  const dst = {};
  for (const item of mdItems) {
    const zoneName = item.name;
    const floor = byZoneName[zoneName];
    if (!floor) continue;
    dst[floor.level] = item;
  }
  return dst;
})
</script>

<style scoped>

</style>