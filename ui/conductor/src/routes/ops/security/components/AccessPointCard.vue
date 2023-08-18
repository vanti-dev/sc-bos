<template>
  <v-card
      color="#40464D"
      elevation="0"
      dark
      min-width="420px"
      height="100%"
      min-height="240px"
      max-height="240px">
    <WithAccess v-slot="{ resource: accessResource }" :name="props.device.name" :paused="props.paused">
      <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <Access
            :access-attempt="accessResource.value"
            :status-log="statusResource.value"
            :loading="accessResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            @click:close="emit('click:close')"/>
      </WithStatus>
    </WithAccess>
  </v-card>
</template>

<script setup>
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
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
  paused: {
    type: Boolean,
    default: false
  },
  showClose: {
    type: Boolean,
    default: false
  }
});
const emit = defineEmits(['click:close']);
</script>
