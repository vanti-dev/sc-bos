<template>
  <v-toolbar color="transparent" elevation="0">
    <v-spacer/>
    <v-select
        hide-details
        filled
        :menu-props="{offsetY: true}"
        label="Floor"
        :items="floors"
        :value="selectedFloor"
        @change="updateFloor"/>
    <v-select
        hide-details
        filled
        :menu-props="{offsetY: true}"
        label="Zone"
        :value="selectedZone"
        :items="zones"
        @change="updateZone"/>
  </v-toolbar>
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  floor: {
    type: String,
    default: ''
  },
  floorItems: {
    type: Array,
    default: () => []
  },
  zone: {
    type: String,
    default: ''
  },
  zoneItems: {
    type: Array,
    default: () => []
  }
});
const emit = defineEmits({
  'update:floor': String,
  'update:zone': String
});

const floors = computed(() => {
  return ['All', ...props.floorItems];
});
const zones = computed(() => {
  return ['All', ...props.zoneItems];
});

const selectedFloor = computed(() => {
  return props.floor || 'All';
});

/**
 *
 * @param {string} f
 */
function updateFloor(f) {
  if (f === 'All') {
    emit('update:floor', undefined);
  } else {
    emit('update:floor', f);
  }
}

const selectedZone = computed(() => {
  return props.zone || 'All';
});

/**
 *
 * @param {string} z
 */
function updateZone(z) {
  if (z === 'All') {
    emit('update:zone', undefined);
  } else {
    emit('update:zone', z);
  }
}
</script>

<style scoped>
::v-deep(.v-toolbar__content > *:not(:last-child)) {
  margin-right: 10px;
}

.v-select {
  max-width: 15em;
}
</style>
