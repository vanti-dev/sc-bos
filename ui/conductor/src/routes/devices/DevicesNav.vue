<template>
  <v-list class="pa-0" dense nav>
    <v-list-item
        v-for="(device, key) in availableSubSystems"
        :id="device.label"
        :key="key"
        :to="device.to"
        @keyup.down="keyboardNavigation($event, device)"
        @keyup.up="keyboardNavigation($event, device)"
        @click="keyboardNavigation($event, device)">
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
import {newActionTracker} from '@/api/resource';
import {useRouter, useRoute} from 'vue-router/composables';

const routeTo = useRouter();
const deviceStore = useDevicesStore();
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
          // to: '/devices/' + subSystem,
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

/**
 * @typedef {Object} Device
 * @property {string} to
 * @property {string} icon
 * @property {string} label
 * @property {string} class
 * @param {KeyboardEvent|MouseEvent} event
 * @param {Device} device
 */
function keyboardNavigation(event, device) {
  const keyCode = event.keyCode;
  const isArrowKey = [38, 40].includes(keyCode); // up and down arrow codes

  if (isArrowKey || event.type === 'click') {
    event.preventDefault();

    const currentIndex = availableSubSystems.value.findIndex(subSystem => subSystem.label === device.label);
    const nextIndex = isArrowKey ? (keyCode === 40 ? currentIndex + 1 : currentIndex - 1) : -1;

    const nextItem = availableSubSystems.value[nextIndex];
    if (nextItem) {
      routeTo.push({path: nextItem.to});
      document.getElementById(nextItem.label).focus();
    } else if (event.type === 'click') {
      routeTo.push({path: device.to});
      event.currentTarget.focus();
    }
  }
}

// onCreate
deviceStore.fetchDeviceSubsystemCounts(tracker);
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: var(--v-primary-base);
}
</style>
