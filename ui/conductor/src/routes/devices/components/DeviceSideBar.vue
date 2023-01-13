<template>
  <side-bar>
    <device-info-card/>
    <air-temperature-card :name="sidebarData.name" v-if="showCard('smartcore.traits.AirTemperature')"/>
  </side-bar>
</template>

<script setup>
import {usePageStore} from '@/stores/page';
import {storeToRefs} from 'pinia';
import AirTemperatureCard from '@/routes/devices/components/trait-cards/AirTemperatureCard.vue';
import DeviceInfoCard from '@/routes/devices/components/trait-cards/DeviceInfoCard.vue';
import SideBar from '@/components/SideBar.vue';
import {ref, watch} from 'vue';

const pageStore = usePageStore();
const {sidebarData} = storeToRefs(pageStore);

const traits = ref([]);
const deviceId = ref('');

watch(sidebarData, (device) => {
  console.log('watch', device);
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
