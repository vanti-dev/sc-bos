<template>
  <span class="root">
    <WithEnterLeave
        v-if="hasCell('EnterLeaveEvent')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <EnterLeaveEventCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else :resource="resource.streamError"/>
    </WithEnterLeave>
    <WithElectricDemand
        v-if="hasCell('ElectricDemand')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <ElectricDemandCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-meter-electric-outline" :resource="resource.streamError"/>
    </WithElectricDemand>
    <WithMeter
        v-if="hasCell('Meter')"
        v-slot="{resource, info}"
        :name="props.item.name"
        :paused="props.paused">
      <MeterCell v-if="!resource.streamError" v-bind="resource" :unit="info?.response?.unit"/>
      <StatusAlert v-else icon="mdi-counter" :resource="resource.streamError"/>
    </WithMeter>
    <WithAirTemperature
        v-if="hasCell('AirTemperature')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <AirTemperatureCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-thermometer-low" :resource="resource.streamError"/>
    </WithAirTemperature>
    <WithLighting
        v-if="hasCell('Light')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <LightCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-lightbulb-outline" :resource="resource.streamError"/>
    </WithLighting>
    <WithOccupancy
        v-if="hasCell('Occupancy')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <OccupancyCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-crosshairs" :resource="resource.streamError"/>
    </WithOccupancy>
    <WithStatus
        v-if="hasCell('StatusLog')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <StatusLogCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-connection" :resource="resource.streamError"/>
    </WithStatus>
    <WithAccess
        v-if="hasCell('AccessAttempt')"
        v-slot="{resource}"
        :name="props.item.name"
        :paused="props.paused">
      <AccessAttemptCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-cancel" :resource="resource.streamError"/>
    </WithAccess>
  </span>
</template>

<script setup>
import WithAirTemperature from '@/routes/devices/components/renderless/WithAirTemperature.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithEnterLeave from '@/routes/devices/components/renderless/WithEnterLeave.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import WithMeter from '@/routes/devices/components/renderless/WithMeter.vue';
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import AirTemperatureCell from '@/routes/devices/components/trait-cells/AirTemperatureCell.vue';
import ElectricDemandCell from '@/routes/devices/components/trait-cells/ElectricDemandCell.vue';
import EnterLeaveEventCell from '@/routes/devices/components/trait-cells/EnterLeaveEventCell.vue';
import StatusLogCell from '@/routes/devices/components/trait-cells/StatusLogCell.vue';
import MeterCell from '@/routes/devices/components/trait-cells/MeterCell.vue';
import AccessAttemptCell from '@/routes/devices/components/trait-cells/AccessAttemptCell.vue';
import {hasTrait} from '@/util/devices';
import {computed} from 'vue';
import WithLighting from './renderless/WithLighting.vue';
import WithOccupancy from './renderless/WithOccupancy.vue';
import LightCell from './trait-cells/LightCell.vue';
import OccupancyCell from './trait-cells/OccupancyCell.vue';
import StatusAlert from '@/components/StatusAlert.vue';

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
const visibleCells = computed(() => {
  const cells = {};
  if (hasTrait(props.item, 'smartcore.traits.OccupancySensor')) {
    cells['Occupancy'] = true;
  }
  if (hasTrait(props.item, 'smartcore.traits.AirTemperature')) {
    cells['AirTemperature'] = true;
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
  if (hasTrait(props.item, 'smartcore.bos.Meter')) {
    cells['Meter'] = true;
  }
  if (hasTrait(props.item, 'smartcore.bos.Access')) {
    cells['AccessAttempt'] = true;
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

