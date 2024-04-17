<template>
  <v-card color="#40464D" elevation="0" dark min-width="420px" height="100%" min-height="240px" max-height="240px">
    <!-- Has Access data but has no OpenClose data -->
    <with-access
        v-if="availableTraits.includes('Access') && !availableTraits.includes('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.device.name"
        :paused="props.paused">
      <with-status v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <access-info
            :access-attempt="accessResource.value"
            :status-log="statusResource.value"
            :loading="accessResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            :paused="props.paused"
            @click:close="emit('click:close')"/>
      </with-status>
    </with-access>
    <!-- Has OpenClose data but has no Access data -->
    <with-open-close
        v-if="!availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: openCloseResource }"
        :name="props.device.name"
        :paused="props.paused">
      <with-status v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
        <access-info
            :open-close="openCloseResource.value"
            :status-log="statusResource.value"
            :loading="openCloseResource.loading || statusResource.loading"
            :device="props.device"
            :show-close="props.showClose"
            :paused="props.paused"
            @click:close="emit('click:close')"/>
      </with-status>
    </with-open-close>
    <!-- Has both Access and OpenClose data -->
    <with-access
        v-if="availableTraits.includes('Access') && availableTraits.includes('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.device.name"
        :paused="props.paused">
      <with-open-close v-slot="{ resource: openCloseResource }" :name="props.device.name" :paused="props.paused">
        <with-status v-slot="{ resource: statusResource }" :name="props.device.name" :paused="props.paused">
          <access-info
              :access-attempt="accessResource.value"
              :open-close="openCloseResource.value"
              :status-log="statusResource.value"
              :loading="accessResource.loading || openCloseResource.loading || statusResource.loading"
              :device="props.device"
              :show-close="props.showClose"
              :paused="props.paused"
              @click:close="emit('click:close')"/>
        </with-status>
      </with-open-close>
    </with-access>
  </v-card>
</template>

<script setup>
import AccessInfo from '@/routes/ops/security/components/access-point-card/AccessInfo.vue';
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
