<template>
  <side-bar>
    <device-info-card/>
    <WithStatus v-if="traits['smartcore.bos.Status']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <status-log-card v-bind="resource"/>
    </WithStatus>
    <WithAirTemperature v-if="traits['smartcore.traits.AirTemperature']" :name="deviceId" v-slot="{resource, update}">
      <v-divider class="mt-4 mb-1"/>
      <air-temperature-card v-bind="resource" @updateAirTemperature="update"/>
    </WithAirTemperature>
    <v-divider v-if="traits['smartcore.traits.Light']" class="mt-4 mb-1"/>
    <WithLighting v-if="traits['smartcore.traits.Light']" :name="deviceId" v-slot="{resource, update}">
      <light-card v-bind="resource" @updateBrightness="update"/>
    </WithLighting>
    <v-divider v-if="traits['smartcore.traits.OccupancySensor']" class="mt-4 mb-1"/>
    <WithOccupancy v-if="traits['smartcore.traits.OccupancySensor']" :name="deviceId" v-slot="{resource}">
      <occupancy-card v-bind="resource"/>
    </WithOccupancy>
    <WithElectricDemand v-if="traits['smartcore.traits.Electric']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <electric-demand-card v-bind="resource"/>
    </WithElectricDemand>
    <WithMeter v-if="traits['smartcore.bos.Meter']" :name="deviceId" v-slot="{resource}">
      <v-divider class="mt-4 mb-1"/>
      <meter-card v-bind="resource"/>
    </WithMeter>
    <v-divider v-if="traits['smartcore.bsp.EmergencyLight']" class="mt-4 mb-1"/>
    <emergency-light :name="deviceId" v-if="traits['smartcore.bsp.EmergencyLight']"/>
    <v-divider v-if="traits['smartcore.traits.Mode']" class="mt-4 mb-1"/>
    <mode-card :name="deviceId" v-if="traits['smartcore.traits.Mode']"/>
    <v-divider v-if="traits['smartcore.bos.UDMI']" class="mt-4 mb-1"/>
    <udmi-card :name="deviceId" v-if="traits['smartcore.bos.UDMI']"/>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import WithAirTemperature from '@/routes/devices/components/renderless/WithAirTemperature.vue';
import WithElectricDemand from '@/routes/devices/components/renderless/WithElectricDemand.vue';
import WithLighting from '@/routes/devices/components/renderless/WithLighting.vue';
import WithOccupancy from '@/routes/devices/components/renderless/WithOccupancy.vue';
import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import WithMeter from '@/routes/devices/components/renderless/WithMeter.vue';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import ElectricDemandCard from '@/routes/devices/components/trait-cards/ElectricDemandCard.vue';
import EmergencyLight from '@/routes/devices/components/trait-cards/EmergencyLight.vue';
import LightCard from '@/routes/devices/components/trait-cards/LightCard.vue';
import ModeCard from '@/routes/devices/components/trait-cards/ModeCard.vue';
import OccupancyCard from '@/routes/devices/components/trait-cards/OccupancyCard.vue';
import StatusLogCard from '@/routes/devices/components/trait-cards/StatusLogCard.vue';
import UdmiCard from '@/routes/devices/components/trait-cards/UdmiCard.vue';
import MeterCard from '@/routes/devices/components/trait-cards/MeterCard.vue';

import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const deviceId = computed(() => {
  return sidebarData.value?.name ?? '';
});

const traits = computed(() => {
  const traits = {};
  if (sidebarData.value?.metadata?.traitsList) {
    // flatten array of trait objects (e.g. [{name: 'trait1', ...}, ...] into object (e.g. {trait1: true, ...})
    sidebarData.value.metadata.traitsList.forEach((trait) => traits[trait.name] = true);
  }
  return traits;
});

</script>

<style scoped>
</style>
