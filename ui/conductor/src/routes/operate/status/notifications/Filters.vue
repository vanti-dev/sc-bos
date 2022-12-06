<template>
  <v-toolbar color="transparent" elevation="0">
    <v-spacer/>
    <v-select hide-details
              filled
              :menu-props="{offsetY: true}"
              label="Floor"
              :items="floors"
              :value="floor"
              @change="updateFloor"></v-select>
    <v-select hide-details
              filled
              :menu-props="{offsetY: true}"
              label="Zone"
              :value="zone"
              :items="zones"
              @change="updateZone"></v-select>
  </v-toolbar>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  floor: String,
  floorItems: Array,
  zone: String,
  zoneItems: Array
});
const emit = defineEmits({
  'update:floor': String,
  'update:zone': String
});

const floors = computed(() => {
  return ['All', ...props.floorItems];
})
const zones = computed(() => {
  return ['All', ...props.zoneItems];
})

const floor = computed(() => {
  return props.floor ?? 'All';
});

function updateFloor(f) {
  if (f === 'All') return emit('update:floor', undefined);
  emit('update:floor', f);
}

const zone = computed(() => {
  return props.zone ?? 'All';
})

function updateZone(z) {
  if (z === 'All') return emit('update:zone', undefined);
  emit('update:zone', z);
}
</script>

<style scoped>
::v-deep(.v-toolbar__content > *:not(:last-child)) {
  margin-right: 10px;
}

.v-select {
  max-width: 120px;
}
</style>
