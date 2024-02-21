<template>
  <span class="root">
    <WithEnterLeave
        v-if="hasCell('EnterLeaveEvent')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <EnterLeaveEventCell
          v-if="resource?.value?.enterTotal || resource?.value?.leaveTotal"
          v-bind="resource"/>
    </WithEnterLeave>
    <WithElectricDemand
        v-if="hasCell('ElectricDemand')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <ElectricDemandCell v-bind="resource"/>
    </WithElectricDemand>
    <WithMeter v-if="hasCell('Meter')" v-slot="{ resource, info }" :name="props.item.name" :paused="props.paused">
      <MeterCell
          v-bind="resource"
          :unit="info?.response?.unit"/>
    </WithMeter>
    <WithAirTemperature
        v-if="hasCell('AirTemperature')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <AirTemperatureCell v-bind="resource"/>
    </WithAirTemperature>

    <LightCell v-if="hasCell('Light')" :name="props.item.name" :paused="props.paused"/>

    <WithOccupancy v-if="hasCell('Occupancy')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <OccupancyCell v-bind="resource"/>
    </WithOccupancy>

    <!-- If door has no access data reading and has OpenClose reading -->
    <WithOpenClosed
        v-if="hasCell('OpenClose') && !hasCell('AccessAttempt')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <OpenClosedCell v-bind="resource"/>
    </WithOpenClosed>

    <!-- If door has access data reading and has no OpenClose reading -->
    <WithAccess
        v-if="hasCell('AccessAttempt') && !hasCell('OpenClose')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <AccessAttemptCell v-bind="resource"/>
    </WithAccess>

    <!-- If door has access data reading and has OpenClose reading -->
    <WithAccess
        v-if="hasCell('AccessAttempt') && hasCell('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.item.name"
        :paused="props.paused">
      <WithOpenClosed v-slot="{ resource: openClosedResource }" :name="props.item.name" :paused="props.paused">
        <AccessAttemptCell
            v-bind="accessResource"
            :open-close-percentage="openClosedResource"
            :stream-error="accessResource.streamError || openClosedResource.streamError"/>
      </WithOpenClosed>
    </WithAccess>
    <!-- End -->

    <WithEmergency v-if="hasCell('Emergency')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <EmergencyCell v-bind="resource"/>
    </WithEmergency>

    <WithStatus v-if="hasCell('StatusLog')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <StatusLogCell v-bind="resource"/>
    </WithStatus>
  </span>
</template>

<script setup>
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import WithAirTemperature from '@/routes/devices/components/renderless/WithAirTemperature.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithEmergency from '@/routes/devices/components/renderless/WithEmergency.vue';
import WithEnterLeave from '@/routes/devices/components/renderless/WithEnterLeave.vue';
import WithMeter from '@/routes/devices/components/renderless/WithMeter.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import WithOpenClosed from '@/routes/devices/components/renderless/WithOpenClosed.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import AccessAttemptCell from '@/traits/access/AccessAttemptCell.vue';
import AirTemperatureCell from '@/traits/airTemperature/AirTemperatureCell.vue';
import ElectricDemandCell from '@/traits/electricDemand/ElectricDemandCell.vue';
import EmergencyCell from '@/traits/emergency/EmergencyCell.vue';
import EnterLeaveEventCell from '@/traits/enterLeave/EnterLeaveEventCell.vue';
import LightCell from '@/traits/lighting/LightCell.vue';
import MeterCell from '@/traits/meter/MeterCell.vue';
import OccupancyCell from '@/traits/occupancy/OccupancyCell.vue';
import OpenClosedCell from '@/traits/openClosed/OpenClosedCell.vue';
import StatusLogCell from '@/traits/status/StatusLogCell.vue';
import {hasTrait} from '@/util/devices';
import {computed} from 'vue';

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
  if (hasTrait(props.item, 'smartcore.traits.Emergency')) {
    cells['Emergency'] = true;
  }
  if (hasTrait(props.item, 'smartcore.traits.OpenClose')) {
    cells['OpenClose'] = true;
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
