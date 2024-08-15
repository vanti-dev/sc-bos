<template>
  <v-toolbar color="transparent" elevation="0">
    <v-spacer/>
    <v-select
        v-for="dropdown in filterDropdowns"
        variant="filled"
        hide-details
        :items="dropdown.items.value"
        :key="dropdown.label"
        :label="dropdown.label"
        :menu-props="{offsetY: true}"
        v-model="dropdown.vModel.value"/>
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
const selectedFloor = computed({
  get() {
    return props.floor || 'All';
  },
  set(f) {
    if (f === 'All') {
      emit('update:floor', undefined);
    } else {
      emit('update:floor', f);
    }
  }
});

const zones = computed(() => {
  return ['All', ...props.zoneItems.filter(v => Boolean(v))];
});
const selectedZone = computed({
  get() {
    return props.zone || 'All';
  },
  set(z) {
    if (z === 'All') {
      emit('update:zone', undefined);
    } else {
      emit('update:zone', z);
    }
  }
});

const subsystems = computed(() => {
  return ['All', ...props.subsystemItems.filter(v => Boolean(v))];
});
const selectedSubsystem = computed({
  get() {
    return props.subsystem || 'All';
  },
  set(s) {
    if (s === 'All') {
      emit('update:subsystem', undefined);
    } else {
      emit('update:subsystem', s);
    }
  }
});

const notificationTypes = computed(() => {
  return ['All', 'Acknowledged', 'Unacknowledged'];
});
const selectedNotificationType = computed({
  get() {
    return props.acknowledged === undefined ? 'All' : props.acknowledged ? 'Acknowledged' : 'Unacknowledged';
  },
  set(a) {
    if (a === 'All') {
      emit('update:acknowledged', undefined);
    } else if (a === 'Acknowledged') {
      emit('update:acknowledged', true);
    } else {
      emit('update:acknowledged', false);
    }
  }
});

const filterDropdowns = [
  {
    label: 'Subsystem',
    items: subsystems,
    vModel: selectedSubsystem
  },
  {
    label: 'Floor',
    items: floors,
    vModel: selectedFloor
  },
  {
    label: 'Zone',
    items: zones,
    vModel: selectedZone
  },
  {
    label: 'Notification Type',
    items: notificationTypes,
    vModel: selectedNotificationType
  }
];
</script>

<style scoped>
:deep(.v-toolbar__content > *:not(:last-child)) {
  margin-right: 10px;
}

.v-select {
  max-width: 15em;
}
</style>
