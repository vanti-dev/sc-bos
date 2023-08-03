<template>
  <v-toolbar color="transparent" elevation="0">
    <v-spacer/>
    <v-select
        hide-details
        filled
        :menu-props="{offsetY: true}"
        label="Subsystem"
        :items="subsystems"
        :value="selectedSubsystem"
        @change="updateSubsystem"/>
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
  },
  subsystem: {
    type: String,
    default: ''
  },
  subsystemItems: {
    type: Array,
    default: () => []
  },
  acknowledged: {
    type: Boolean,
    default: undefined
  },
  resolved: {
    type: Boolean,
    default: undefined
  }
});
const emit = defineEmits({
  'update:floor': String,
  'update:zone': String,
  'update:subsystem': String,
  'update:acknowledged': Boolean,
  'update:resolved': Boolean
});

const floors = computed(() => {
  return ['All', ...props.floorItems.filter(v => Boolean(v))];
});
const zones = computed(() => {
  return ['All', ...props.zoneItems.filter(v => Boolean(v))];
});
const subsystems = computed(() => {
  return ['All', ...props.subsystemItems.filter(v => Boolean(v))];
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

const selectedSubsystem = computed(() => {
  return props.subsystem || 'All';
});

/**
 *
 * @param {string} s
 */
function updateSubsystem(s) {
  if (s === 'All') {
    emit('update:subsystem', undefined);
  } else {
    emit('update:subsystem', s);
  }
}
</script>

<style scoped>
:deep(.v-toolbar__content > *:not(:last-child)) {
  margin-right: 10px;
}

.v-select {
  max-width: 15em;
}
</style>
