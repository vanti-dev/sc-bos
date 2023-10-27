<template>
  <span class="root">
    <WithEnterLeave
        v-if="hasCell('EnterLeaveEvent')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <EnterLeaveEventCell
          v-if="!resource.streamError && (resource?.value?.enterTotal || resource?.value?.leaveTotal)"
          v-bind="resource"/>
      <StatusAlert v-else :resource="resource.streamError"/>
    </WithEnterLeave>
    <WithElectricDemand
        v-if="hasCell('ElectricDemand')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <ElectricDemandCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-meter-electric-outline" :resource="resource.streamError"/>
    </WithElectricDemand>
    <WithMeter v-if="hasCell('Meter')" v-slot="{ resource, info }" :name="props.item.name" :paused="props.paused">
      <MeterCell v-if="!resource.streamError" v-bind="resource" :unit="info?.response?.unit"/>
      <StatusAlert v-else icon="mdi-counter" :resource="resource.streamError"/>
    </WithMeter>
    <WithAirTemperature
        v-if="hasCell('AirTemperature')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <AirTemperatureCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-thermometer-low" :resource="resource.streamError"/>
    </WithAirTemperature>
    <WithLighting v-if="hasCell('Light')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <LightCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-lightbulb-outline" :resource="resource.streamError"/>
    </WithLighting>
    <WithOccupancy v-if="hasCell('Occupancy')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <OccupancyCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-crosshairs" :resource="resource.streamError"/>
    </WithOccupancy>

    <!-- If door has no access data reading and has OpenClose reading -->
    <WithOpenClosed
        v-if="hasCell('OpenClose') && !hasCell('AccessAttempt')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <OpenClosedCell v-if="!resource?.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-cancel" :resource="resource.streamError"/>
    </WithOpenClosed>

    <!-- If door has access data reading and has no OpenClose reading -->
    <WithAccess
        v-if="hasCell('AccessAttempt') && !hasCell('OpenClose')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <AccessAttemptCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-cancel" :resource="resource.streamError"/>
    </WithAccess>

    <!-- If door has access data reading and has OpenClose reading -->
    <WithAccess
        v-if="hasCell('AccessAttempt') && hasCell('OpenClose')"
        v-slot="{ resource: accessResource }"
        :name="props.item.name"
        :paused="props.paused">
      <WithOpenClosed v-slot="{ resource: openClosedResource }" :name="props.item.name" :paused="props.paused">
        <AccessAttemptCell
            v-if="!accessResource.streamError || !openClosedResource.streamError"
            v-bind="accessResource"
            :open-close-percentage="openClosedResource"/>
        <StatusAlert
            v-else-if="accessResource.streamError"
            icon="mdi-cancel"
            :resource="accessResource.streamError"/>
        <StatusAlert
            v-else-if="openClosedResource.streamError"
            icon="mdi-cancel"
            :resource="openClosedResource.streamError"/>
      </WithOpenClosed>
    </WithAccess>
    <!-- End -->

    <WithEmergency v-if="hasCell('Emergency')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <EmergencyCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-smoke-detector-outline" :resource="resource.streamError"/>
    </WithEmergency>

    <WithStatus v-if="hasCell('StatusLog')" v-slot="{ resource }" :name="props.item.name" :paused="props.paused">
      <StatusLogCell v-if="!resource.streamError" v-bind="resource"/>
      <StatusAlert v-else icon="mdi-connection" :resource="resource.streamError"/>
    </WithStatus>
  </span>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import WithAccess from '@/routes/devices/components/renderless/WithAccess.vue';
import WithAirTemperature from '@/routes/devices/components/renderless/WithAirTemperature.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithEmergency from '@/routes/devices/components/renderless/WithEmergency.vue';
import WithEnterLeave from '@/routes/devices/components/renderless/WithEnterLeave.vue';
import WithLighting from '@/routes/devices/components/renderless/WithLighting.vue';
import WithMeter from '@/routes/devices/components/renderless/WithMeter.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import WithOpenClosed from '@/routes/devices/components/renderless/WithOpenClosed.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import AccessAttemptCell from '@/routes/devices/components/trait-cells/AccessAttemptCell.vue';
import AirTemperatureCell from '@/routes/devices/components/trait-cells/AirTemperatureCell.vue';
import ElectricDemandCell from '@/routes/devices/components/trait-cells/ElectricDemandCell.vue';
import EmergencyCell from '@/routes/devices/components/trait-cells/EmergencyCell.vue';
import EnterLeaveEventCell from '@/routes/devices/components/trait-cells/EnterLeaveEventCell.vue';
import LightCell from '@/routes/devices/components/trait-cells/LightCell.vue';
import MeterCell from '@/routes/devices/components/trait-cells/MeterCell.vue';
import OccupancyCell from '@/routes/devices/components/trait-cells/OccupancyCell.vue';
import OpenClosedCell from '@/routes/devices/components/trait-cells/OpenClosedCell.vue';
import StatusLogCell from '@/routes/devices/components/trait-cells/StatusLogCell.vue';
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
