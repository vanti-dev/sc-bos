<template>
  <v-list class="pa-0" density="compact" nav>
    <v-list-item
        v-for="(device, key) in availableSubSystems"
        :key="key"
        :to="device.to"
        :disabled="hasNoAccess(device.to)"
        class="my-2">
      <template #prepend>
        <v-icon v-if="device.icon">{{ device.icon }}</v-icon>
        <subsystem-icon v-else :subsystem="device.subSystem"/>
      </template>
      <v-list-item-title :class="[device.class, 'text-truncate']">
        {{ device.label }}
      </v-list-item-title>
    </v-list-item>
  </v-list>
</template>

<script setup>
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import {useDevicesMetadataField, usePullDevicesMetadata} from '@/composables/devices.js';
import useAuthSetup from '@/composables/useAuthSetup';
import {computed} from 'vue';

const {hasNoAccess} = useAuthSetup();
const {value: md} = usePullDevicesMetadata('metadata.membership.subsystem');
const {keys} = useDevicesMetadataField(md, 'metadata.membership.subsystem');

// computed
//
// Generating the Devices nav options
const availableSubSystems = computed(() => {
  // Pinia store value
  const subSystems = keys.value;

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
    for (const subSystem of subSystems) {
      // ignore noType key/value pair
      if (subSystem !== '') {
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
</script>

<style scoped>
:deep(.v-list-item--active) {
  color: rgb(var(--v-theme-primary));
}
</style>
