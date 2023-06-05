<template>
  <span class="root">
    <WithEnterLeave
        v-if="hasCell('EnterLeaveEvent')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <EnterLeaveEventCell v-bind="resource"/>
    </WithEnterLeave>
    <WithElectricDemand
        v-if="hasCell('ElectricDemand')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <ElectricDemandCell v-bind="resource"/>
    </WithElectricDemand>
    <WithLighting
        v-if="hasCell('Light')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <LightCell v-bind="resource"/>
    </WithLighting>
    <WithOccupancy
        v-if="hasCell('Occupancy')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <OccupancyCell v-bind="resource"/>
    </WithOccupancy>
    <WithStatus
        v-if="hasCell('StatusLog')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <StatusLogCell v-bind="resource"/>
    </WithStatus>
  </span>
</template>

<script setup>
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithEnterLeave from '@/routes/devices/components/renderless/WithEnterLeave.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import ElectricDemandCell from '@/routes/devices/components/trait-cells/ElectricDemandCell.vue';
import EnterLeaveEventCell from '@/routes/devices/components/trait-cells/EnterLeaveEventCell.vue';
import StatusLogCell from '@/routes/devices/components/trait-cells/StatusLogCell.vue';
import {hasTrait} from '@/util/devices';
import {computed} from 'vue';
import WithLighting from './renderless/WithLighting.vue';
import WithOccupancy from './renderless/WithOccupancy.vue';
import LightCell from './trait-cells/LightCell.vue';
import OccupancyCell from './trait-cells/OccupancyCell.vue';

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
const visibleCells = computed(() => {
  const cells = {};
  if (hasTrait(props.item, 'smartcore.traits.OccupancySensor')) {
    cells['Occupancy'] = true;
  }
  if (hasTrait(props.item, 'smartcore.traits.Light')) {
    cells['Light'] = true;
  }
  if (hasTrait(props.item, 'smartcore.traits.Electric')) {
    cells['ElectricDemand'] = true;
  }
  if (hasTrait(props.item, 'smartcore.traits.EnterLeaveSensor')) {
    cells['EnterLeaveEvent'] = true;
  }
  if (hasTrait(props.item, 'smartcore.bos.Status')) {
    cells['StatusLog'] = true;
  }
  return cells;
});

/**
 * @param {string} name
 * @return {boolean}
 */
function hasCell(name) {
  return Boolean(visibleCells.value[name]);
}
</script>

<style scoped>
.root {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-end;
  gap: 1em;
}
</style>
