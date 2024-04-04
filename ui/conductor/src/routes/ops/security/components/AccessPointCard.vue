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
    <WithOpenClose
        v-if="!availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: openCloseResource }"
        :name="props.device.name"
        :paused="props.paused">
      <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <Access
            :open-close="openCloseResource.value"
            :status-log="statusResource.value"
            :loading="openCloseResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            :paused="props.paused"
            @click:close="emit('click:close')"/>
      </WithStatus>
    </WithOpenClose>
    <!-- Has both Access and OpenClose data -->
    <WithAccess
        v-if="availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.device.name"
        :paused="props.paused">
      <WithOpenClose v-slot="{ resource: openCloseResource }" :name="props.device.name" :paused="props.paused">
        <WithStatus v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
          <Access
              :access-attempt="accessResource.value"
              :open-close="openCloseResource.value"
              :status-log="statusResource.value"
              :loading="accessResource.loading || openCloseResource.loading || statusResource.loading"
              :device="props.device"
              :show-close="props.showClose"
              :paused="props.paused"
              @click:close="emit('click:close')"/>
        </WithStatus>
      </WithOpenClose>
    </WithAccess>
  </v-card>
</template>

<script setup>
import Access from '@/routes/ops/security/components/access-point-card/Access.vue';
import WithAccess from '@/traits/access/WithAccess.vue';
import WithOpenClose from '@/traits/openClose/WithOpenClose.vue';
import WithStatus from '@/traits/status/WithStatus.vue';
import {computed} from 'vue';

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
