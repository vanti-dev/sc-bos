<template>
  <v-container fluid class="d-flex flex-row flex-wrap justify-space-between pt-10 mt-14 mx-0 px-3">
    <WithAccess
        v-for="(device, deviceIndex) in props.deviceNames"
        :key="deviceIndex"
        :name="device.source"
        class="mb-8"
        v-slot="{ resource }">
      <AccessPointCard v-bind="resource" :source="device.source" :name="device.name"/>
    </WithAccess>
  </v-container>
</template>

<script setup>
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import AccessPointCard from '@/routes/ops/security/components/AccessPointCard.vue';

const props = defineProps({
  deviceNames: {
    type: Array,
    default: () => []
  },
  filter: {
    type: Function,
    default: () => ({})
  }
});
</script>
