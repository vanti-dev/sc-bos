<template>
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
</template>

<script setup>
import WithOccupancy from '../renderless/WithOccupancy.vue';
import OccupancyCell from './OccupancyCell.vue';
import WithLighting from '../renderless/WithLighting.vue';
import LightCell from './LightCell.vue';

import {hasTrait} from '@/util/devices';

const props = defineProps({
  paused: {
    type: Boolean,
    default: false
  },
  item: {
    type: Object,
    default: () => {}
  }
});
</script>
