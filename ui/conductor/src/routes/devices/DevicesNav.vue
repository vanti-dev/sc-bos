<template>
  <v-list class="pa-0" dense nav>
    <v-list-item
        v-for="(device, key) in availableSubSystems"
        :key="key"
        :to="device.to"
        @click="resetIntersectedItemNames()">
      <v-list-item-icon>
        <v-icon>{{ device.icon }}</v-icon>
      </v-list-item-icon>
      <v-list-item-content :class="device.class">{{ device.label }}</v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {computed, reactive} from 'vue';
import {useDevicesStore} from './store';
import {useTableDataStore} from '@/stores/tableDataStore';
import {newActionTracker} from '@/api/resource';

const deviceStore = useDevicesStore();
const {resetIntersectedItemNames} = useTableDataStore();
const tracker = reactive(newActionTracker());

// computed
//
// Generating the Devices nav options
const availableSubSystems = computed(() => {
  // Pinia store value
  const subSystems = deviceStore.subSystems.subs;

  // navigationItemIcons
  const navigationItemIcons = {
    lighting: 'mdi-lightbulb',
    hvac: 'mdi-thermometer',
    metering: 'mdi-meter-electric',
    acs: 'mdi-badge-account-horizontal',
    cctv: 'mdi-cctv',
    fire: 'mdi-fire',
    smart: 'mdi-memory',
    vt: 'mdi-elevator-passenger',
    sensor: 'mdi-leak', // we might need to change this icon
    zones: 'mdi-select-all'
  };

  const navigationItemLabels = {
    acs: 'Access Control',
    vt: 'Vertical Transportation'
  };


  // array of navigationItems with a default item value - 'All'
  const navigationItems = [
    {to: '/devices/all', icon: 'mdi-view-list', label: 'all', class: 'text-capitalize'}
  ];

  // if we have the Pinia store value
  if (subSystems) {
    // then loop through this value
    for (const subSystem in subSystems) {
      // ignore noType key/value pair
      if (subSystem !== 'noType') {
        // and reconstruct the object according to the template needs
        const listItem = {
          to: '/devices/' + encodeURIComponent(subSystem),
          icon: navigationItemIcons[subSystem] ? navigationItemIcons[subSystem] : 'mdi-chevron-right',
          label: navigationItemLabels[subSystem] ? navigationItemLabels[subSystem] : subSystem,
          class: ['hvac', 'cctv'].includes(subSystem) ? 'text-uppercase' : 'text-capitalize'
        };

        // then fill the array with the new object(s)
        navigationItems.push(listItem);
      }
    }
  }

  // finally return the array of objects
  return navigationItems;
});

// onCreate
deviceStore.fetchDeviceSubsystemCounts(tracker);
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
