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
          :info="info?.response"/>
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
    <WithOpenClose
        v-if="hasCell('OpenClose') && !hasCell('AccessAttempt')"
        v-slot="{ resource }"
        :name="props.item.name"
        :paused="props.paused">
      <OpenCloseCell v-bind="resource"/>
    </WithOpenClose>

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
      <WithOpenClose v-slot="{ resource: openCloseResource }" :name="props.item.name" :paused="props.paused">
        <AccessAttemptCell
            v-bind="accessResource"
            :open-close-percentage="openCloseResource"
            :stream-error="accessResource.streamError || openCloseResource.streamError"/>
      </WithOpenClose>
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
import AccessAttemptCell from '@/traits/access/AccessAttemptCell.vue';
import WithAccess from '@/traits/access/WithAccess.vue';
import AirTemperatureCell from '@/traits/airTemperature/AirTemperatureCell.vue';
import WithAirTemperature from '@/traits/airTemperature/WithAirTemperature.vue';
import ElectricDemandCell from '@/traits/electricDemand/ElectricDemandCell.vue';
import WithElectricDemand from '@/traits/electricDemand/WithElectricDemand.vue';
import EmergencyCell from '@/traits/emergency/EmergencyCell.vue';
import WithEmergency from '@/traits/emergency/WithEmergency.vue';
import EnterLeaveEventCell from '@/traits/enterLeave/EnterLeaveEventCell.vue';
import WithEnterLeave from '@/traits/enterLeave/WithEnterLeave.vue';
import LightCell from '@/traits/light/LightCell.vue';
import MeterCell from '@/traits/meter/MeterCell.vue';
import WithMeter from '@/traits/meter/WithMeter.vue';
import OccupancyCell from '@/traits/occupancy/OccupancyCell.vue';
import WithOccupancy from '@/traits/occupancy/WithOccupancy.vue';
import OpenCloseCell from '@/traits/openClose/OpenCloseCell.vue';
import WithOpenClose from '@/traits/openClose/WithOpenClose.vue';
import StatusLogCell from '@/traits/status/StatusLogCell.vue';
import WithStatus from '@/traits/status/WithStatus.vue';
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
