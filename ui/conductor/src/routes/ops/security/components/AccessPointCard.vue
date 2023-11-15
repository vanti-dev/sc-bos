<template>
  <v-card color="#40464D" elevation="0" dark min-width="420px" height="100%" min-height="240px" max-height="240px">
    <!-- Has Access data but has no OpenClose data -->
    <WithAccess
        v-if="availableTraits.includes('Access') && !availableTraits.includes('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.device.name"
        :paused="props.paused">
      <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <Access
            :access-attempt="accessResource.value"
            :status-log="statusResource.value"
            :loading="accessResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            :paused="props.paused"
            @click:close="emit('click:close')"/>
      </WithStatus>
    </WithAccess>
    <!-- Has OpenClose data but has no Access data -->
    <WithOpenClosed
        v-if="!availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: openClosedResource }"
        :name="props.device.name"
        :paused="props.paused">
      <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <Access
            :open-closed="openClosedResource.value"
            :status-log="statusResource.value"
            :loading="openClosedResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            :paused="props.paused"
            @click:close="emit('click:close')"/>
      </WithStatus>
    </WithOpenClosed>
    <!-- Has both Access and OpenClose data -->
    <WithAccess
        v-if="availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.device.name"
        :paused="props.paused">
      <WithOpenClosed v-slot="{ resource: openClosedResource }" :name="props.device.name" :paused="props.paused">
        <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
          <Access
              :access-attempt="accessResource.value"
              :open-closed="openClosedResource.value"
              :status-log="statusResource.value"
              :loading="accessResource.loading || openClosedResource.loading || statusResource.loading"
              :device="props.device"
              :show-close="props.showClose"
              :paused="props.paused"
              @click:close="emit('click:close')"/>
        </WithStatus>
      </WithOpenClosed>
    </WithAccess>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';

import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import WithOpenClosed from '@/routes/devices/components/renderless/WithOpenClosed.vue';
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

const availableTraits = computed(() => {
  return props.device.traits.map((trait) => trait.split('.').at(-1)) || [];
});
</script>
