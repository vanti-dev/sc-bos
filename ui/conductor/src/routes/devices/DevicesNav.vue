<template>
  <v-list class="pa-0" dense nav>
    <v-list-item
        v-for="(device, key) in availableSubSystems"
        :key="key"
        :to="device.to"
        :disabled="hasNoAccess(device.to)"
        class="my-2">
      <v-list-item-icon>
        <v-icon v-if="device.icon">{{ device.icon }}</v-icon>
        <subsystem-icon v-else :subsystem="device.subSystem"/>
      </v-list-item-icon>
      <v-list-item-content :class="[device.class, 'text-truncate']">
        {{ device.label }}
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script setup>
import {newActionTracker} from '@/api/resource';
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import useAuthSetup from '@/composables/useAuthSetup';
import {computed, reactive} from 'vue';
import {useDevicesStore} from './store';

const {hasNoAccess} = useAuthSetup();
const deviceStore = useDevicesStore();
const tracker = reactive(newActionTracker());

// computed
//
// Generating the Devices nav options
const availableSubSystems = computed(() => {
  // Pinia store value
  const subSystems = deviceStore.subSystems.subs;

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
          subSystem,
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
