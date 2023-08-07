<template>
  <v-container fluid class="d-flex flex-row flex-wrap justify-space-between pt-4 mx-0 px-3">
    <WithAccess
        v-for="(device, deviceIndex) in deviceNames"
        :key="deviceIndex"
        :name="device.source"
        class="mb-8"
        v-slot="{ resource }">
      <AccessPointCard v-bind="resource" :source="device.source" :name="device.name"/>
    </WithAccess>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import AccessPointCard from '@/routes/ops/security/components/AccessPointCard.vue';

const props = defineProps({
  devices: {
    type: Array,
    default: () => []
  },
  filter: {
    type: Function,
    default: () => ({})
  }
});

const deviceNames = computed(() => {
  return props.devices.map((device) => {
    return {
      source: device.metadata.name,
      name: device.metadata?.appearance ? device.metadata?.appearance.title : device.metadata.name
    };
  }
  );
});
</script>
