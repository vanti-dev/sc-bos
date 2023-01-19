<template>
  <side-bar>
    <device-info-card/>
    <air-temperature-card :name="deviceId" v-if="traits['smartcore.traits.AirTemperature']"/>
    <light-card :name="deviceId" v-if="traits['smartcore.traits.Light']"/>
    <emergency-light :name="deviceId" v-if="traits['smartcore.bsp.EmergencyLight']"/>
  </side-bar>
</template>

<script setup>
import {computed, watch} from 'vue';
import {storeToRefs} from 'pinia';

import {usePageStore} from '@/stores/page';

import SideBar from '@/components/SideBar.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import LightCard from '@/routes/devices/components/trait-cards/LightCard.vue';
import EmergencyLight from '@/routes/devices/components/trait-cards/EmergencyLight.vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const deviceId = computed(() => {
  return sidebarData.value?.name ?? '';
});

const traits = computed(() => {
  const traits = {};
  if (sidebarData.value?.metadata?.traitsList !== undefined) {
    // flatten array of trait objects (e.g. [{name: 'trait1', ...}, ...] into object (e.g. {trait1: true, ...})
    return sidebarData.value.metadata.traitsList.reduce((obj, key) => ({...obj, [key.name]: true}), {});
  }
  return traits;
});

</script>

<style scoped>
</style>
