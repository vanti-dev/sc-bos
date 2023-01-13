<template>
  <side-bar>
    <device-info-card/>
    <air-temperature-card :name="sidebarData.name" v-if="showCard('smartcore.traits.AirTemperature')"/>
    <light-card :name="sidebarData.name" v-if="showCard('smartcore.traits.Light')"/>
  </side-bar>
</template>

<script setup>
import {ref, watch} from 'vue';
import {storeToRefs} from 'pinia';

import {usePageStore} from '@/stores/page';

import SideBar from '@/components/SideBar.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import LightCard from '@/routes/devices/components/trait-cards/LightCard.vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const traits = ref([]);
const deviceId = ref('');

watch(sidebarData, (device) => {
  deviceId.value = device.name;
  traits.value = [];
  if (device &&
      device.hasOwnProperty('metadata') &&
      device.metadata.hasOwnProperty('traitsList')) {
    traits.value = device.metadata.traitsList.map(t => t.name);
  }
});

/**
 * @param {string} trait
 * @return {boolean}
 */
function showCard(trait) {
  return traits.value.indexOf(trait) >= 0;
}

</script>

<style scoped>
.sidebarTitle {
  background: var(--v-neutral-lighten1);
  height: auto;
  font-weight: bold;
}
</style>
