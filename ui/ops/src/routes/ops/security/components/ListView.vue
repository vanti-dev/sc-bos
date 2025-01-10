<template>
  <v-container fluid class="d-flex flex-row flex-wrap justify-left pt-14 mt-14 mx-0 px-2">
    <span v-for="(device) in props.deviceNames" :key="device.name" class="mb-8 mx-auto">
      <hot-point v-slot="{ live }" :item-key="device.name">
        <access-point-card :device="device" :paused="!live"/>
      </hot-point>
    </span>
  </v-container>
</template>

<script setup>
import HotPoint from '@/components/HotPoint.vue';
import AccessPointCard from '@/routes/ops/security/components/AccessPointCard.vue';

const props = defineProps({
  deviceNames: {
    type: Array,
    default: () => []
  }
});
</script>
