<template>
  <side-bar>
    <device-info-card/>
    <air-temperature-card :name="deviceId" v-if="traits['smartcore.traits.AirTemperature']"/>
    <light-card :name="deviceId" v-if="traits['smartcore.traits.Light']"/>
    <occupancy-card :name="deviceId" v-if="traits['smartcore.traits.OccupancySensor']"/>
    <emergency-light :name="deviceId" v-if="traits['smartcore.bsp.EmergencyLight']"/>
    <mode-card :name="deviceId" v-if="traits['smartcore.traits.Mode']"/>
    <udmi-card :name="deviceId" v-if="traits['smartcore.bos.UDMI']"/>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import EmergencyLight from '@/routes/devices/components/trait-cards/EmergencyLight.vue';
import LightCard from '@/routes/devices/components/trait-cards/LightCard.vue';
import ModeCard from '@/routes/devices/components/trait-cards/ModeCard.vue';
import OccupancyCard from '@/routes/devices/components/trait-cards/OccupancyCard.vue';
import UdmiCard from '@/routes/devices/components/trait-cards/UdmiCard.vue';

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
