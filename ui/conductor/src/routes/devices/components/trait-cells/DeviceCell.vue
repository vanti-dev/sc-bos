<template>
  <span>
    <WithOccupancy
        v-if="hasTrait(props.item, 'OccupancySensor')"
        v-slot="{value}"
        :name="props.item.name"
        :paused="props.paused">
      <OccupancyCell v-bind="{value}"/>
    </WithOccupancy>
    <WithLighting
        v-else-if="hasTrait(props.item, 'Light')"
        v-slot="value"
        :name="props.item.name"
        :paused="props.paused">
      <LightCell v-bind="value"/>
    </WithLighting>
    <WithEnterLeave
        v-else-if="hasTrait(props.item, 'smartcore.traits.EnterLeaveSensor')"
        v-slot="value"
        :name="props.item.name"
        :paused="props.paused">
      <EnterLeaveEventCell v-bind="value"/>
    </WithEnterLeave>
    <WithStatus
        v-if="hasTrait(props.item, 'smartcore.bos.Status')"
        v-slot="value"
        :name="props.item.name"
        :paused="props.paused">
      <StatusLogCell v-bind="value"/>
    </WithStatus>
  </span>
</template>

<script setup>
import WithEnterLeave from '@/routes/devices/components/renderless/WithEnterLeave.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import EnterLeaveEventCell from '@/routes/devices/components/trait-cells/EnterLeaveEventCell.vue';
import StatusLogCell from '@/routes/devices/components/trait-cells/StatusLogCell.vue';
import {hasTrait} from '@/util/devices';
import WithLighting from '../renderless/WithLighting.vue';
import WithOccupancy from '../renderless/WithOccupancy.vue';
import LightCell from './LightCell.vue';
import OccupancyCell from './OccupancyCell.vue';

const props = defineProps({
  paused: {
    type: Boolean,
    default: false
  },
  item: {
    type: Object,
    default: () => {
    }
  }
});
</script>
