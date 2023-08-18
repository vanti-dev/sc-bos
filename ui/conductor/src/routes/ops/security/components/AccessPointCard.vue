<template>
  <v-card
      color="#40464D"
      elevation="0"
      dark
      min-width="420px"
      height="100%"
      min-height="240px"
      max-height="240px">
    <WithAccess v-slot="{ resource }" :name="props.device.name" :paused="props.paused">
      <Access v-bind="resource" :device="props.device"/>
    </WithAccess>
  </v-card>
</template>

<script setup>
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import Access from '@/routes/ops/security/components/access-point-card/Access.vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => ({})
  },
  device: {
    type: Object,
    default: () => ({})
  },
  floor: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  }
});
</script>
