<template>
  <v-container fluid class="d-flex flex-row flex-wrap justify-space-between pt-4 mx-0 px-3">
    <AccessPointCard
        v-for="(deviceName, deviceIndex) in deviceNames"
        :key="deviceIndex"
        class="mb-8"
        :name="deviceName"/>
  </v-container>
</template>

<script setup>
import {computed} from 'vue';
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
  return props.devices.map((device) => device.metadata.name);
});
</script>
